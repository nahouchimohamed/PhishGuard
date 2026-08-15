package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/phishguard/phishguard/internal/store"
)

func TestMapStoreErr(t *testing.T) {
	cases := []struct {
		err  error
		code int
		msg  string
	}{
		{store.ErrNotFound, http.StatusNotFound, "not found"},
		{store.ErrUnauthorized, http.StatusUnauthorized, "unauthorized"},
		{store.ErrForbidden, http.StatusForbidden, "forbidden"},
		{store.ErrInvalidInput, http.StatusBadRequest, "invalid input"},
		{store.ErrConflict, http.StatusConflict, "conflict"},
		{errors.New("boom"), http.StatusInternalServerError, "internal error"},
	}
	for _, tt := range cases {
		rec := httptest.NewRecorder()
		mapStoreErr(rec, tt.err)
		if rec.Code != tt.code {
			t.Fatalf("%v: status %d want %d", tt.err, rec.Code, tt.code)
		}
		var body map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["error"] != tt.msg {
			t.Fatalf("%v: error=%q want %q", tt.err, body["error"], tt.msg)
		}
	}
}

func TestPublicUserJSONKeys(t *testing.T) {
	u := store.User{
		ID:          "u1",
		OrgID:       "o1",
		Email:       "a@example.com",
		DisplayName: "Ada",
		Role:        store.RoleAdmin,
	}
	raw, err := json.Marshal(publicUser(u))
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	want := []string{"id", "org_id", "email", "display_name", "role"}
	if len(m) != len(want) {
		t.Fatalf("keys=%v want exactly %v", m, want)
	}
	for _, k := range want {
		if _, ok := m[k]; !ok {
			t.Fatalf("missing key %q in %v", k, m)
		}
	}
	if m["id"] != "u1" || m["org_id"] != "o1" || m["email"] != "a@example.com" {
		t.Fatalf("unexpected values %#v", m)
	}
	if m["display_name"] != "Ada" || m["role"] != string(store.RoleAdmin) {
		t.Fatalf("unexpected values %#v", m)
	}
}
