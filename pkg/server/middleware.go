package server

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

// corsMiddleware returns a middleware which sets CORS headers on all responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", viper.GetString("cors-allow-origin"))
		next.ServeHTTP(w, r)
	})
}

// SecurityHeadersHandler returns a middleware which sets common security
// HTTP headers on the response to mitigate common web vulnerabilities.
func SecurityHeadersHandler(next http.Handler) http.Handler {
	csp := []string{
		"default-src 'self'",
		"font-src 'self' data:",
		"form-action 'self'",
		"frame-ancestors 'none'",
		"img-src 'self' data:",
		"script-src 'self'",
		"style-src 'self' 'unsafe-inline'",
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-security-policy", strings.Join(csp, "; "))
		w.Header().Set("referrer-policy", "no-referrer")
		w.Header().Set("x-content-type-options", "nosniff")
		w.Header().Set("x-frame-options", "DENY")
		w.Header().Set("x-xss-protection", "1; mode=block")
		if r.URL.Scheme == "https" || r.Header.Get("X-Forwarded-Proto") == "https" {
			w.Header().Set("strict-transport-security", "max-age=31536000")
		}
		next.ServeHTTP(w, r)
	})
}

// newMetricsMiddleware creates a middleware handler recording all HTTP requests in
// the given Prometheus registry
func newMetricsMiddleware(reg prometheus.Registerer) func(http.Handler) http.Handler {
	requests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "yopass_http_requests_total",
			Help: "Total number of requests served by HTTP method, path and response code.",
		},
		[]string{"method", "path", "code"},
	)
	reg.MustRegister(requests)

	duration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "yopass_http_request_duration_seconds",
			Help:    "Histogram of HTTP request latencies by method and path.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
	reg.MustRegister(duration)

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := statusCodeRecorder{ResponseWriter: w, statusCode: http.StatusOK}
			handler.ServeHTTP(&rec, r)
			path := normalizedPath(r)
			requests.WithLabelValues(r.Method, path, strconv.Itoa(rec.statusCode)).Inc()
			duration.WithLabelValues(r.Method, path).Observe(time.Since(start).Seconds())
		})
	}
}

// normalizedPath returns a normalized mux path template representation
func normalizedPath(r *http.Request) string {
	if route := mux.CurrentRoute(r); route != nil {
		if tmpl, err := route.GetPathTemplate(); err == nil {
			return strings.ReplaceAll(tmpl, keyParameter, ":key")
		}
	}
	return "<other>"
}

// statusCodeRecorder is a HTTP ResponseWriter recording the response code
type statusCodeRecorder struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader implements http.ResponseWriter
func (rw *statusCodeRecorder) WriteHeader(code int) {
	rw.ResponseWriter.WriteHeader(code)
	rw.statusCode = code
}
