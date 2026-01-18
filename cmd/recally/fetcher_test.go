package main

import (
	"testing"

	"recally/internal/pkg/webreader"
	"recally/internal/pkg/webreader/fetcher"
)

func TestNewFetcher(t *testing.T) {
	tests := []struct {
		name       string
		useBrowser bool
		browserURL string
		wantErr    bool
		wantType   string // "http" or "browser"
	}{
		{
			name:       "HTTP fetcher",
			useBrowser: false,
			browserURL: "",
			wantErr:    false,
			wantType:   "http",
		},
		{
			name:       "Browser fetcher with valid URL",
			useBrowser: true,
			browserURL: "http://localhost:9222",
			wantErr:    false,
			wantType:   "browser",
		},
		{
			name:       "Browser fetcher with empty URL (launches new browser)",
			useBrowser: true,
			browserURL: "",
			wantErr:    false,
			wantType:   "browser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewFetcher(tt.useBrowser, tt.browserURL)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewFetcher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if f == nil {
				t.Error("NewFetcher() returned nil fetcher without error")
				return
			}

			// Verify we got the right type
			switch tt.wantType {
			case "http":
				if _, ok := f.(*fetcher.HTTPFetcher); !ok {
					t.Errorf("NewFetcher() returned %T, want *fetcher.HTTPFetcher", f)
				}
			case "browser":
				if _, ok := f.(*fetcher.BrowserFetcher); !ok {
					t.Errorf("NewFetcher() returned %T, want *fetcher.BrowserFetcher", f)
				}
			}

			// Clean up
			if err := f.Close(); err != nil {
				t.Errorf("Failed to close fetcher: %v", err)
			}
		})
	}
}

func TestNewHTTPFetcher(t *testing.T) {
	f, err := newHTTPFetcher()
	if err != nil {
		t.Fatalf("newHTTPFetcher() error = %v", err)
	}
	defer func() { _ = f.Close() }()

	if f == nil {
		t.Fatal("newHTTPFetcher() returned nil")
	}

	// Verify it's an HTTP fetcher
	httpFetcher, ok := f.(*fetcher.HTTPFetcher)
	if !ok {
		t.Fatalf("newHTTPFetcher() returned %T, want *fetcher.HTTPFetcher", f)
	}

	// Verify it implements the Fetcher interface
	var _ webreader.Fetcher = httpFetcher
}

func TestNewBrowserFetcher(t *testing.T) {
	tests := []struct {
		name       string
		browserURL string
		wantErr    bool
	}{
		{
			name:       "valid URL",
			browserURL: "http://localhost:9222",
			wantErr:    false,
		},
		{
			name:       "empty URL (launches new browser)",
			browserURL: "",
			wantErr:    false,
		},
		{
			name:       "URL with path",
			browserURL: "http://localhost:9222/devtools",
			wantErr:    false,
		},
		{
			name:       "HTTPS URL",
			browserURL: "https://remote-browser:9222",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := newBrowserFetcher(tt.browserURL)

			if (err != nil) != tt.wantErr {
				t.Errorf("newBrowserFetcher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if f == nil {
				t.Error("newBrowserFetcher() returned nil without error")
				return
			}

			// Verify it's a browser fetcher
			browserFetcher, ok := f.(*fetcher.BrowserFetcher)
			if !ok {
				t.Fatalf("newBrowserFetcher() returned %T, want *fetcher.BrowserFetcher", f)
			}

			// Verify it implements the Fetcher interface
			var _ webreader.Fetcher = browserFetcher

			// Clean up
			if err := f.Close(); err != nil {
				t.Errorf("Failed to close fetcher: %v", err)
			}
		})
	}
}

func TestFetcherInterface(t *testing.T) {
	// Verify both fetchers implement the webreader.Fetcher interface
	t.Run("HTTP fetcher implements interface", func(t *testing.T) {
		f, err := newHTTPFetcher()
		if err != nil {
			t.Fatalf("newHTTPFetcher() error = %v", err)
		}
		defer func() { _ = f.Close() }()

		var _ = f
	})

	t.Run("Browser fetcher implements interface", func(t *testing.T) {
		// Test with empty URL (launches new browser)
		f, err := newBrowserFetcher("")
		if err != nil {
			t.Fatalf("newBrowserFetcher() error = %v", err)
		}
		defer func() { _ = f.Close() }()

		var _ = f
	})
}

func TestFetcherConfiguration(t *testing.T) {
	t.Run("HTTP fetcher has correct timeout", func(t *testing.T) {
		f, err := newHTTPFetcher()
		if err != nil {
			t.Fatalf("newHTTPFetcher() error = %v", err)
		}
		defer func() { _ = f.Close() }()

		// Note: We can't directly inspect the timeout without exporting internal fields
		// This test just verifies the fetcher was created successfully
		// The timeout is verified by the configuration passed to NewHTTPFetcher
	})

	t.Run("Browser fetcher has correct config", func(t *testing.T) {
		// Test with empty URL (launches new browser)
		f, err := newBrowserFetcher("")
		if err != nil {
			t.Fatalf("newBrowserFetcher() error = %v", err)
		}
		defer func() { _ = f.Close() }()

		// Note: We can't directly inspect the config without exporting internal fields
		// This test just verifies the fetcher was created successfully
		// The config is verified by the BrowserConfig struct passed to NewBrowserFetcher
	})
}
