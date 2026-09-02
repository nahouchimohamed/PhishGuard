package metrics

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func counterValue(t *testing.T, name string, labels map[string]string) float64 {
	t.Helper()
	mfs, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatal(err)
	}
	for _, mf := range mfs {
		if mf.GetName() != name {
			continue
		}
		for _, m := range mf.GetMetric() {
			if labelsMatch(m.GetLabel(), labels) {
				if m.Counter != nil {
					return m.GetCounter().GetValue()
				}
			}
		}
	}
	return 0
}

func labelsMatch(got []*dto.LabelPair, want map[string]string) bool {
	if len(want) == 0 {
		return true
	}
	have := map[string]string{}
	for _, lp := range got {
		have[lp.GetName()] = lp.GetValue()
	}
	for k, v := range want {
		if have[k] != v {
			return false
		}
	}
	return true
}

func TestMiddlewareRecordsStatus(t *testing.T) {
	before := counterValue(t, "phishguard_http_requests_total", map[string]string{
		"method": "GET", "path": "/metrics-test", "status": "204",
	})
	h := Middleware("/metrics-test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status=%d", rec.Code)
	}
	after := counterValue(t, "phishguard_http_requests_total", map[string]string{
		"method": "GET", "path": "/metrics-test", "status": "204",
	})
	if after-before != 1 {
		t.Fatalf("delta=%v want 1", after-before)
	}
}

func TestMiddlewareRecordsOnPanicPathViaDefer(t *testing.T) {
	before := counterValue(t, "phishguard_http_requests_total", map[string]string{
		"method": "GET", "path": "/metrics-panic", "status": "200",
	})
	h := Middleware("/metrics-panic", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "ok")
		panic("boom")
	}))
	rec := httptest.NewRecorder()
	func() {
		defer func() { _ = recover() }()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
	}()
	after := counterValue(t, "phishguard_http_requests_total", map[string]string{
		"method": "GET", "path": "/metrics-panic", "status": "200",
	})
	if after-before != 1 {
		t.Fatalf("defer should still record: delta=%v", after-before)
	}
}
