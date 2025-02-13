package processor

import (
	"fmt"
	"regexp"
	"strings"
)

func parseXmlContent(content, tag string) string {
	re := regexp.MustCompile(fmt.Sprintf(`(?s)<%s>(.*?)</%s>`, tag, tag))
	matches := re.FindStringSubmatch(content)
	if len(matches) < 2 {
		return ""
	}
	return strings.TrimSpace(matches[1])
}

func tagStringToArray(tagString string) []string {
	tags := []string{}
	for _, tag := range strings.Split(tagString, ", ") {
		if tag != "" {
			tags = append(tags, strings.TrimSpace(tag))
		}
	}
	return tags
}
