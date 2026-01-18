package main

import (
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid http url",
			url:     "http://example.com/article",
			wantErr: false,
		},
		{
			name:    "valid https url",
			url:     "https://example.com/article",
			wantErr: false,
		},
		{
			name:    "empty url",
			url:     "",
			wantErr: true,
		},
		{
			name:    "url without scheme",
			url:     "example.com",
			wantErr: true,
		},
		{
			name:    "file scheme",
			url:     "file:///etc/passwd",
			wantErr: true,
		},
		{
			name:    "javascript scheme",
			url:     "javascript:alert(1)",
			wantErr: true,
		},
		{
			name:    "ftp scheme",
			url:     "ftp://example.com",
			wantErr: true,
		},
		{
			name:    "data scheme",
			url:     "data:text/html,<html></html>",
			wantErr: true,
		},
		{
			name:    "url without host",
			url:     "https://",
			wantErr: true,
		},
		{
			name:    "valid url with path and query",
			url:     "https://example.com/path?key=value",
			wantErr: false,
		},
		{
			name:    "valid url with fragment",
			url:     "https://example.com/article#section",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
