package bookmarks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseXmlContent(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		tag         string
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:    "successfully parse content with tag",
			content: "<title>Test Title</title>",
			tag:     "title",
			want:    "Test Title",
			wantErr: false,
		},
		{
			name:    "successfully parse content with multiline tag",
			content: "<description>\nThis is a\nmultiline\ndescription\n</description>",
			tag:     "description",
			want:    "\nThis is a\nmultiline\ndescription\n",
			wantErr: false,
		},
		{
			name:        "missing tag returns error",
			content:     "<other>Some content</other>",
			tag:         "title",
			want:        "",
			wantErr:     true,
			errContains: "no title tag found in content",
		},
		{
			name:        "empty content returns error",
			content:     "",
			tag:         "title",
			want:        "",
			wantErr:     true,
			errContains: "no title tag found in content",
		},
		{
			name:    "handle content with multiple tags",
			content: "<title>First Title</title><description>Some desc</description><title>Second Title</title>",
			tag:     "title",
			want:    "First Title",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseXmlContent(tt.content, tt.tag)
			assert.Equal(t, tt.want, got)
		})
	}
}
