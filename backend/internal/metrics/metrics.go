package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	CampaignLaunches = promauto.NewCounter(prometheus.CounterOpts{
		Name: "phishguard_campaign_launches_total",
		Help: "Campaigns successfully moved to launching/active",
	})
	Deliveries = promauto.NewCounter(prometheus.CounterOpts{
		Name: "phishguard_deliveries_total",
		Help: "Successful inbox deliveries by the worker",
	})
	Clicks = promauto.NewCounter(prometheus.CounterOpts{
		Name: "phishguard_clicks_total",
		Help: "Simulation link clicks recorded",
	})
	Reports = promauto.NewCounter(prometheus.CounterOpts{
		Name: "phishguard_reports_total",
		Help: "Employee report actions",
	})
	HTTPRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "phishguard_http_requests_total",
		Help: "HTTP requests by method, path template, and status",
	}, []string{"method", "path", "status"})
	HTTPDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "phishguard_http_request_duration_seconds",
		Help:    "HTTP request latency",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})
)

func Handler() http.Handler {
	return promhttp.Handler()
}

type statusWriter struct {
	http.ResponseWriter
	code int
}

func (w *statusWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Middleware records latency and status for one path label (keep cardinality low).
func Middleware(pathLabel string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, code: http.StatusOK}
		next.ServeHTTP(sw, r)
		HTTPRequests.WithLabelValues(r.Method, pathLabel, strconv.Itoa(sw.code)).Inc()
		HTTPDuration.WithLabelValues(r.Method, pathLabel).Observe(time.Since(start).Seconds())
	})
}
