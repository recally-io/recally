package bookmarks

import (
	"testing"
)

func TestParseTagsFromSummary(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedTags []string
		expectedText string
	}{
		{
			name:         "Normal tags",
			input:        "Hello World\n<tags>\n#tag1 #tag2 #tag3\n</tags>\nGoodbye",
			expectedTags: []string{"tag1", "tag2", "tag3"},
			expectedText: "Hello World\n\n\nGoodbye",
		},
		{
			name:         "No tags section",
			input:        "Hello World",
			expectedTags: []string{},
			expectedText: "Hello World",
		},
		{
			name:         "Empty tags section",
			input:        "Hello World<tags></tags>Goodbye",
			expectedTags: []string{},
			expectedText: "Hello World\nGoodbye",
		},
		{
			name:         "Invalid tags",
			input:        "Hello World<tags>invalid# # invalid</tags>Goodbye",
			expectedTags: []string{},
			expectedText: "Hello World\nGoodbye",
		},
		{
			name:         "Duplicate tags",
			input:        "Hello<tags>#tag1 #tag2 #tag1</tags>World",
			expectedTags: []string{"tag1", "tag2"},
			expectedText: "Hello\nWorld",
		},
		{
			name:         "Mixed valid and invalid tags",
			input:        "Start<tags>#valid1 invalid# #valid2 # #valid3</tags>End",
			expectedTags: []string{"valid1", "valid2", "valid3"},
			expectedText: "Start\nEnd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags, text := parseTagsFromSummary(tt.input)

			if !equalSlices(tags, tt.expectedTags) {
				t.Errorf("tags = %v, want %v", tags, tt.expectedTags)
			}

			if text != tt.expectedText {
				t.Errorf("text = %v, want %v", text, tt.expectedText)
			}
		})
	}
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
