package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// SecretURL returns a URL which decrypts the specified secret in the browser.
func SecretURL(baseURL, id, key string, fileOpt, manualKeyOpt bool) string {
	prefix := "s"
	if fileOpt {
		prefix = "f"
	}
	path := id
	if !manualKeyOpt {
		path += "/" + key
	}
	return fmt.Sprintf("%s/#/%s/%s", strings.TrimSuffix(baseURL, "/"), prefix, path)
}

// ParseURL returns secret ID and key from a regular yopass URL.
func ParseURL(s string) (id, key string, fileOpt, keyOpt bool, err error) {
	u, err := url.Parse(strings.TrimSpace(s))
	if err != nil {
		return "", "", false, false, fmt.Errorf("invalid URL: %w", err)
	}

	f := strings.Split(u.Fragment, "/")
	if len(f) < 3 || len(f) > 4 || f[0] != "" {
		return "", "", false, false, fmt.Errorf("unexpected URL: %q", s)
	}

	switch f[1] {
	case "s":
	case "c":
		keyOpt = true
	case "f":
		fileOpt = true
	case "d":
		fileOpt = true
		keyOpt = true
	default:
		return "", "", false, false, fmt.Errorf("unexpected URL: %q", s)
	}

	id = f[2]
	if len(f) == 4 {
		key = f[3]
	}
	return id, key, fileOpt, keyOpt, nil
}
