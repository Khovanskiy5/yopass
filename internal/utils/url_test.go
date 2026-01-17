package utils

import (
	"testing"
)

func TestSecretURL(t *testing.T) {
	tests := []struct {
		name         string
		baseURL      string
		id           string
		key          string
		fileOpt      bool
		manualKeyOpt bool
		expected     string
	}{
		{
			name:         "Secret URL",
			baseURL:      "https://yopass.se",
			id:           "uuid",
			key:          "secretkey",
			fileOpt:      false,
			manualKeyOpt: false,
			expected:     "https://yopass.se/#/s/uuid/secretkey",
		},
		{
			name:         "File URL",
			baseURL:      "https://yopass.se/",
			id:           "uuid",
			key:          "secretkey",
			fileOpt:      true,
			manualKeyOpt: false,
			expected:     "https://yopass.se/#/f/uuid/secretkey",
		},
		{
			name:         "Manual key option",
			baseURL:      "https://yopass.se",
			id:           "uuid",
			key:          "secretkey",
			fileOpt:      false,
			manualKeyOpt: true,
			expected:     "https://yopass.se/#/s/uuid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SecretURL(tt.baseURL, tt.id, tt.key, tt.fileOpt, tt.manualKeyOpt)
			if got != tt.expected {
				t.Errorf("SecretURL() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseURL(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		wantId       string
		wantKey      string
		wantFileOpt  bool
		wantKeyOpt   bool
		wantErr      bool
	}{
		{
			name:        "Valid secret URL",
			url:         "https://yopass.se/#/s/uuid/key",
			wantId:      "uuid",
			wantKey:     "key",
			wantFileOpt: false,
			wantKeyOpt:  false,
			wantErr:     false,
		},
		{
			name:        "Valid file URL",
			url:         "https://yopass.se/#/f/uuid/key",
			wantId:      "uuid",
			wantKey:     "key",
			wantFileOpt: true,
			wantKeyOpt:  false,
			wantErr:     false,
		},
		{
			name:        "Valid custom key URL",
			url:         "https://yopass.se/#/c/uuid",
			wantId:      "uuid",
			wantKey:     "",
			wantFileOpt: false,
			wantKeyOpt:  true,
			wantErr:     false,
		},
		{
			name:        "Valid download URL",
			url:         "https://yopass.se/#/d/uuid",
			wantId:      "uuid",
			wantKey:     "",
			wantFileOpt: true,
			wantKeyOpt:  true,
			wantErr:     false,
		},
		{
			name:    "Invalid URL",
			url:     "https://yopass.se/s/uuid/key",
			wantErr: true,
		},
		{
			name:    "Unexpected prefix",
			url:     "https://yopass.se/#/x/uuid/key",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, key, fileOpt, keyOpt, err := ParseURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if id != tt.wantId || key != tt.wantKey || fileOpt != tt.wantFileOpt || keyOpt != tt.wantKeyOpt {
				t.Errorf("ParseURL() = %v, %v, %v, %v; want %v, %v, %v, %v", id, key, fileOpt, keyOpt, tt.wantId, tt.wantKey, tt.wantFileOpt, tt.wantKeyOpt)
			}
		})
	}
}
