package bookmarks

import (
	"regexp"
	"strings"
)

func parseListFilter(filters []string) (domains, contentTypes, tags []string) {
	if len(filters) == 0 {
		return
	}

	domains = make([]string, 0)
	contentTypes = make([]string, 0)
	tags = make([]string, 0)

	// Parse filter=category:article;type:rss
	for _, part := range filters {
		kv := strings.Split(part, ":")
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "domain":
			domains = append(domains, kv[1])
		case "type":
			contentTypes = append(contentTypes, kv[1])
		case "tag":
			tags = append(tags, kv[1])
		}
	}
	if len(domains) == 0 {
		domains = nil
	}
	if len(contentTypes) == 0 {
		contentTypes = nil
	}
	if len(tags) == 0 {
		tags = nil
	}
	return
}

// parseTagsFromSummary extracts tags from a string and returns the tags array and the string without tags
func parseTagsFromSummary(input string) ([]string, string) {
	// Regular expression to match the tags section
	tagsRegex := regexp.MustCompile(`(?s)<tags>.*?</tags>`)

	// Find tags section
	tagsSection := tagsRegex.FindString(input)

	// If no tags section found, return empty array and original string
	if tagsSection == "" {
		return []string{}, input
	}

	// Extract content between tags
	content := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(tagsSection, "<tags>"), "</tags>"))

	// Split content by whitespace
	words := strings.Fields(content)

	// Process valid tags
	tagMap := make(map[string]bool) // Use map to ensure uniqueness
	var tags []string

	for _, word := range words {
		if strings.HasPrefix(word, "#") {
			tag := strings.TrimPrefix(word, "#")
			if tag != "" && !tagMap[tag] {
				tagMap[tag] = true
				tags = append(tags, tag)
			}
		}
	}

	// Remove tags section from original string
	cleanedString := strings.TrimSpace(tagsRegex.ReplaceAllString(input, "\n"))

	return tags, cleanedString
}
