package utils

import (
	"net/http"
	"testing"
)

func TestGetRealClientIP(t *testing.T) {
	tests := []struct {
		name           string
		remoteAddr     string
		xForwardedFor  string
		trustedProxies []string
		expected       string
	}{
		{
			name:       "No trusted proxies, simple remote addr",
			remoteAddr: "1.2.3.4:1234",
			expected:   "1.2.3.4",
		},
		{
			name:           "Trusted proxy CIDR match",
			remoteAddr:     "10.0.0.1:1234",
			xForwardedFor:  "1.2.3.4, 10.0.0.1",
			trustedProxies: []string{"10.0.0.0/24"},
			expected:       "1.2.3.4",
		},
		{
			name:           "Trusted proxy IP match",
			remoteAddr:     "10.0.0.1:1234",
			xForwardedFor:  "1.2.3.4",
			trustedProxies: []string{"10.0.0.1"},
			expected:       "1.2.3.4",
		},
		{
			name:           "Untrusted proxy, use remote addr",
			remoteAddr:     "10.0.0.1:1234",
			xForwardedFor:  "1.2.3.4",
			trustedProxies: []string{"192.168.1.1"},
			expected:       "10.0.0.1",
		},
		{
			name:           "Trusted proxy, invalid X-Forwarded-For",
			remoteAddr:     "10.0.0.1:1234",
			xForwardedFor:  "invalid-ip",
			trustedProxies: []string{"10.0.0.1"},
			expected:       "10.0.0.1",
		},
		{
			name:           "Trusted proxy, empty X-Forwarded-For",
			remoteAddr:     "10.0.0.1:1234",
			xForwardedFor:  "",
			trustedProxies: []string{"10.0.0.1"},
			expected:       "10.0.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			req.RemoteAddr = tt.remoteAddr
			if tt.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.xForwardedFor)
			}

			got := GetRealClientIP(req, tt.trustedProxies)
			if got != tt.expected {
				t.Errorf("GetRealClientIP() = %v, want %v", got, tt.expected)
			}
		})
	}
}
