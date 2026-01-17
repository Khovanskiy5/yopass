package utils

import (
	"net"
	"net/http"
	"strings"
)

// GetRealClientIP returns the real client IP address by checking X-Forwarded-For
// header only if the request comes from a trusted proxy, otherwise returns RemoteAddr
func GetRealClientIP(req *http.Request, trustedProxies []string) string {
	remoteAddr := req.RemoteAddr

	// Extract IP from RemoteAddr (removes port if present)
	remoteIP, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		remoteIP = remoteAddr
	}

	// If no trusted proxies configured, always use RemoteAddr
	if len(trustedProxies) == 0 {
		return remoteIP
	}

	// Check if the request comes from a trusted proxy
	isTrusted := false
	for _, trustedProxy := range trustedProxies {
		// Parse CIDR or single IP
		_, cidr, err := net.ParseCIDR(trustedProxy)
		if err != nil {
			// Not a CIDR, try as single IP
			if remoteIP == trustedProxy {
				isTrusted = true
				break
			}
		} else {
			// Check if remoteIP is in the CIDR range
			if cidr.Contains(net.ParseIP(remoteIP)) {
				isTrusted = true
				break
			}
		}
	}

	// If not from trusted proxy, use RemoteAddr to prevent spoofing
	if !isTrusted {
		return remoteIP
	}

	// Extract the first IP from X-Forwarded-For header
	xForwardedFor := req.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs separated by commas
		// The first IP is the original client IP
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			clientIP := strings.TrimSpace(ips[0])
			if net.ParseIP(clientIP) != nil {
				return clientIP
			}
		}
	}

	// Fallback to RemoteAddr if X-Forwarded-For is invalid or empty
	return remoteIP
}
