package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestCORS(t *testing.T) {
	allowOrigin := "https://example.com"
	m := CORS(allowOrigin)
	
	handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "http://localhost", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != allowOrigin {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, allowOrigin)
	}
}

func TestSecurityHeaders(t *testing.T) {
	handler := SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name   string
		scheme string
		header string
		exists bool
	}{
		{"CSP", "http", "Content-Security-Policy", true},
		{"Referrer-Policy", "http", "Referrer-Policy", true},
		{"X-Content-Type-Options", "http", "X-Content-Type-Options", true},
		{"X-Frame-Options", "http", "X-Frame-Options", true},
		{"X-XSS-Protection", "http", "X-XSS-Protection", true},
		{"HSTS HTTP", "http", "Strict-Transport-Security", false},
		{"HSTS HTTPS", "https", "Strict-Transport-Security", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://localhost", nil)
			if tt.scheme == "https" {
				req.URL.Scheme = "https"
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			val := w.Header().Get(tt.header)
			if tt.exists && val == "" {
				t.Errorf("Header %s should be present", tt.header)
			}
			if !tt.exists && val != "" {
				t.Errorf("Header %s should not be present", tt.header)
			}
		})
	}
}

func TestMetrics(t *testing.T) {
	reg := prometheus.NewRegistry()
	m := Metrics(reg)

	handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest("POST", "/api/secret", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	// Verify that metrics were recorded
	metricFamilies, err := reg.Gather()
	if err != nil {
		t.Fatalf("Gather failed: %v", err)
	}

	found := false
	for _, mf := range metricFamilies {
		if mf.GetName() == "yopass_http_requests_total" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Metric yopass_http_requests_total not found")
	}
}
