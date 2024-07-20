package handlers // replace with your actual package name

import (
	"testing"
)

// Test_getUrlFromText tests the getUrlFromText function for various scenarios.
func Test_getUrlFromText(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "Test with valid URL",
			text: "Check out this link: https://example.com",
			want: "https://example.com",
		},
		{
			name: "Test with valid http URL",
			text: "Check out this link: http://example.com",
			want: "http://example.com",
		},
		{
			name: "Test with text but no URL",
			text: "Just some random text without a URL",
			want: "",
		},
		{
			name: "Test with URL PATH and params",
			text: "Check out this link: https://example.com/path/to/something?param1=value1&param2=value2 xxx",
			want: "https://example.com/path/to/something?param1=value1&param2=value2",
		},
		{
			name: "Test with multiple URLs",
			text: "Here are two URLs: https://example.com and https://another-example.com",
			want: "https://example.com", // Assuming the function returns the first URL found
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getUrlFromText(tt.text); got != tt.want {
				t.Errorf("getUrlFromText() = %v, want %v", got, tt.want)
			}
		})
	}
}
