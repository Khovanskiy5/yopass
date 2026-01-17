package constants

import (
	"regexp"
	"testing"
)

func TestKeyParameterRegex(t *testing.T) {
	// KeyParameter is a gorilla/mux path variable with a regex: {key:regex}
	// We need to extract the regex part to test it
	re := regexp.MustCompile(`\{key:(.*)\}`)
	matches := re.FindStringSubmatch(KeyParameter)
	if len(matches) < 2 {
		t.Fatalf("Failed to extract regex from KeyParameter: %s", KeyParameter)
	}
	pattern := matches[1]
	keyRegex, err := regexp.Compile("^" + pattern + "$")
	if err != nil {
		t.Fatalf("Failed to compile extracted regex %s: %v", pattern, err)
	}

	tests := []struct {
		key   string
		valid bool
	}{
		{"21701b28-fb3f-451d-8a52-3e6c9094e7ea", true},
		{"00000000-0000-0000-0000-000000000000", true},
		{"ffffffff-ffff-ffff-ffff-ffffffffffff", true},
		{"21701b28-fb3f-451d-8a52-3e6c9094e7e", false},  // too short
		{"21701b28-fb3f-451d-8a52-3e6c9094e7eaa", false}, // too long
		{"21701b28fb3f451d8a523e6c9094e7ea", false},     // missing dashes
		{"g1701b28-fb3f-451d-8a52-3e6c9094e7ea", false}, // invalid hex
	}

	for _, tc := range tests {
		t.Run(tc.key, func(t *testing.T) {
			if keyRegex.MatchString(tc.key) != tc.valid {
				t.Errorf("Expected match=%v for key %s", tc.valid, tc.key)
			}
		})
	}
}
