package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/phishguard/phishguard/internal/auth"
	"github.com/phishguard/phishguard/internal/campaign"
	"github.com/phishguard/phishguard/internal/metrics"
	"github.com/phishguard/phishguard/internal/store"
)

type Server struct {
	auth      *auth.Service
	store     *store.Store
	campaigns *campaign.Service
	mux       *http.ServeMux
}

func New(authSvc *auth.Service, st *store.Store, camps *campaign.Service) *Server {
	s := &Server{auth: authSvc, store: st, campaigns: camps, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) routes() {
	s.mux.Handle("GET /healthz", metrics.Middleware("/healthz", http.HandlerFunc(s.handleHealth)))
	s.mux.Handle("GET /metrics", metrics.Handler())

	s.mux.Handle("POST /api/login", metrics.Middleware("/api/login", http.HandlerFunc(s.handleLogin)))
	s.mux.Handle("POST /api/logout", metrics.Middleware("/api/logout", http.HandlerFunc(s.handleLogout)))
	s.mux.Handle("GET /api/me", metrics.Middleware("/api/me", s.requireAuth(s.handleMe)))

	s.mux.Handle("GET /api/employees", metrics.Middleware("/api/employees", s.requireRole(store.RoleAdmin, s.handleListEmployees)))
	s.mux.Handle("GET /api/templates", metrics.Middleware("/api/templates", s.requireRole(store.RoleAdmin, s.handleListTemplates)))
	s.mux.Handle("GET /api/campaigns", metrics.Middleware("/api/campaigns", s.requireRole(store.RoleAdmin, s.handleListCampaigns)))
	s.mux.Handle("POST /api/campaigns", metrics.Middleware("/api/campaigns", s.requireRole(store.RoleAdmin, s.handleCreateCampaign)))
	s.mux.Handle("GET /api/campaigns/{id}", metrics.Middleware("/api/campaigns/{id}", s.requireRole(store.RoleAdmin, s.handleGetCampaign)))
	s.mux.Handle("POST /api/campaigns/{id}/retry", metrics.Middleware("/api/campaigns/{id}/retry", s.requireRole(store.RoleAdmin, s.handleRetryCampaign)))

	s.mux.Handle("GET /api/inbox", metrics.Middleware("/api/inbox", s.requireRole(store.RoleEmployee, s.handleInbox)))
	s.mux.Handle("GET /api/inbox/{id}", metrics.Middleware("/api/inbox/{id}", s.requireRole(store.RoleEmployee, s.handleInboxItem)))
	s.mux.Handle("POST /api/inbox/{id}/open", metrics.Middleware("/api/inbox/{id}/open", s.requireRole(store.RoleEmployee, s.handleOpen)))
	s.mux.Handle("POST /api/inbox/{id}/report", metrics.Middleware("/api/inbox/{id}/report", s.requireRole(store.RoleEmployee, s.handleReport)))

	s.mux.Handle("GET /api/sim/{token}", metrics.Middleware("/api/sim/{token}", http.HandlerFunc(s.handleSim)))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := s.store.Ping(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, "db unavailable")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	user, token, err := s.auth.Login(r.Context(), strings.TrimSpace(strings.ToLower(req.Email)), req.Password)
	if err != nil {
		mapStoreErr(w, err)
		return
	}
	s.auth.SetSessionCookie(w, token)
	// session_token is for the SvelteKit BFF to set the browser cookie when it
	// proxies login. Direct API clients can rely on Set-Cookie alone.
	writeJSON(w, http.StatusOK, map[string]any{
		"user":          publicUser(user),
		"session_token": token,
	})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie(s.auth.CookieName())
	token := ""
	if c != nil {
		token = c.Value
	}
	_ = s.auth.Logout(r.Context(), token)
	s.auth.ClearSessionCookie(w)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	user := userFromCtx(r.Context())
	writeJSON(w, http.StatusOK, publicUser(user))
}

func (s *Server) handleListEmployees(w http.ResponseWriter, r *http.Request) {
	user := userFromCtx(r.Context())
	emps, err := s.store.ListEmployees(r.Context(), user.OrgID)
	if err != nil {
		mapStoreErr(w, err)
		return
	}
	out := make([]map[string]any, 0, len(emps))
	for _, e := range emps {
		out = append(out, publicUser(e))
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) handleListTemplates(w http.ResponseWriter, r *http.Request) {
	user := userFromCtx(r.Context())
	tmpls, err := s.store.ListTemplates(r.Context(), user.OrgID)
	if err != nil {
		mapStoreErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tmpls)
}

func (s *Server) handleListCampaigns(w http.ResponseWriter, r *http.Request) {
	user := userFromCtx(r.Context())
	camps, err := s.store.ListCampaigns(r.Context(), user.OrgID)
	if err != nil {
		mapStoreErr(w, err)
		return
	}
	type row struct {
		store.Campaign
		Stats store.CampaignStats `json:"stats"`
	}
	out := make([]row, 0, len(camps))
	for _, c := range camps {
		st, err := s.campaigns.GetStats(r.Context(), c.ID)
		if err != nil {
			mapStoreErr(w, err)
			return
		}
		out = append(out, row{Campaign: c, Stats: st})
	}
	writeJSON(w, http.StatusOK, out)
}

type createCampaignRequest struct {
	Name             string   `json:"name"`
	TemplateID       string   `json:"template_id"`
	RecipientUserIDs []string `json:"recipient_user_ids"`
}

func (s *Server) handleCreateCampaign(w http.ResponseWriter, r *http.Request) {
	user := userFromCtx(r.Context())
	var req createCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	camp, err := s.campaigns.Launch(r.Context(), campaign.LaunchInput{
		OrgID:            user.OrgID,
		CreatedBy:        user.ID,
		Name:             strings.TrimSpace(req.Name),
		TemplateID:       req.TemplateID,
		RecipientUserIDs: req.RecipientUserIDs,
	})
	if err != nil {
		// Campaign may still exist in launching state; return 202-ish info via error message for the UI.
		if camp.ID != "" {
			writeJSON(w, http.StatusAccepted, map[string]any{
				"campaign": camp,
				"warning":  err.Error(),
			})
			return
		}
		mapStoreErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, camp)
}

func (s *Server) handleGetCampaign(w http.ResponseWriter, r *http.Request) {
	user := userFromCtx(r.Context())
	id := r.PathValue("id")
	camp, err := s.store.GetCampaign(r.Context(), id, user.OrgID)
	if err != nil {
		mapStoreErr(w, err)
		return
	}
	tmpl, err := s.store.GetTemplate(r.Context(), camp.TemplateID, user.OrgID)
	if err != nil {
		mapStoreErr(w, err)
		return
	}
	stats, err := s.campaigns.GetStats(r.Context(), camp.ID)
	if err != nil {
		mapStoreErr(w, err)
		return
	}
	recs, err := s.store.ListRecipientRows(r.Context(), camp.ID)
	if err != nil {
		mapStoreErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"campaign":   camp,
		"template":   tmpl,
		"stats":      stats,
		"recipients": recs,
	})
}

func (s *Server) handleRetryCampaign(w http.ResponseWriter, r *http.Request) {
	user := userFromCtx(r.Context())
	id := r.PathValue("id")
	camp, err := s.campaigns.RetryLaunch(r.Context(), id, user.OrgID)
	if err != nil {
		mapStoreErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, camp)
}

func (s *Server) handleInbox(w http.ResponseWriter, r *http.Request) {
	user := userFromCtx(r.Context())
	items, err := s.store.ListInbox(r.Context(), user.ID)
	if err != nil {
		mapStoreErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleInboxItem(w http.ResponseWriter, r *http.Request) {
	user := userFromCtx(r.Context())
	id := r.PathValue("id")
	item, err := s.store.GetInboxItem(r.Context(), id, user.ID)
	if err != nil {
		mapStoreErr(w, err)
		return
	}
	_ = s.campaigns.Open(r.Context(), id, user.ID, user.OrgID)
	item.OpenedAt = ptrTime(time.Now().UTC())
	item.SimURL = s.campaigns.SimURL(item.ClickToken)
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) handleOpen(w http.ResponseWriter, r *http.Request) {
	user := userFromCtx(r.Context())
	id := r.PathValue("id")
	if err := s.campaigns.Open(r.Context(), id, user.ID, user.OrgID); err != nil {
		mapStoreErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleReport(w http.ResponseWriter, r *http.Request) {
	user := userFromCtx(r.Context())
	id := r.PathValue("id")
	if err := s.campaigns.Report(r.Context(), id, user.ID, user.OrgID); err != nil {
		mapStoreErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "reported"})
}

func (s *Server) handleSim(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	tmpl, campaignID, err := s.campaigns.Click(r.Context(), token)
	if err != nil {
		mapStoreErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"campaign_id": campaignID,
		"subject":     tmpl.Subject,
		"lesson_text": tmpl.LessonText,
		"name":        tmpl.Name,
	})
}

func ptrTime(t time.Time) *time.Time { return &t }

func publicUser(u store.User) map[string]any {
	return map[string]any{
		"id":           u.ID,
		"org_id":       u.OrgID,
		"email":        u.Email,
		"display_name": u.DisplayName,
		"role":         u.Role,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func mapStoreErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, "not found")
	case errors.Is(err, store.ErrUnauthorized):
		writeError(w, http.StatusUnauthorized, "unauthorized")
	case errors.Is(err, store.ErrForbidden):
		writeError(w, http.StatusForbidden, "forbidden")
	case errors.Is(err, store.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, "invalid input")
	case errors.Is(err, store.ErrConflict):
		writeError(w, http.StatusConflict, "conflict")
	default:
		slog.Error("request failed", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}
