package processor

import (
	"context"
	"fmt"
	"net/url"
	"recally/internal/pkg/webreader"
	"strings"

	"github.com/go-shiori/go-readability"
)

type ReadabilityProcessor struct{}

func NewReadabilityProcessor() *ReadabilityProcessor {
	return &ReadabilityProcessor{}
}

func (p *ReadabilityProcessor) Name() string {
	return "Readability"
}

// Process implements the Processor interface
func (p *ReadabilityProcessor) Process(ctx context.Context, content *webreader.Content) error {
	parsedURL, err := url.ParseRequestURI(content.URL)
	// If there's an error parsing the URI, set parsedURL to nil
	if err != nil {
		parsedURL = nil
	}

	// run before hooks
	// beforeHooks := hooks.GetReadabilityBeforeHooks(parsedURL.Host)
	// if len(beforeHooks) > 0 {
	// 	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(content.Html))
	// 	for _, hook := range beforeHooks {
	// 		hook(doc, content)
	// 	}
	// }

	// Use the readability package's FromReader function to parse the HTML content
	article, err := readability.FromReader(strings.NewReader(content.Html), parsedURL)
	// If there's an error parsing the HTML content, return the error
	if err != nil {
		return fmt.Errorf("failed to parse %s, %v", content.URL, err)
	}

	// run after hooks
	// afterHooks := hooks.GetReadabilityAfterHooks(parsedURL.Host)
	// if len(afterHooks) > 0 {
	// 	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(content.Html))
	// 	for _, hook := range afterHooks {
	// 		hook(doc, content)
	// 	}
	// }

	// Set meta info
	content.Title = article.Title
	content.Author = article.Byline
	content.Description = article.Excerpt
	content.SiteName = article.SiteName

	content.Cover = article.Image
	content.Favicon = article.Favicon

	// set text content
	content.Text = article.TextContent
	// set clean HTML content processed by readability
	content.Html = article.Content

	// Set the published and modified time
	content.PublishedTime = article.PublishedTime
	content.ModifiedTime = article.ModifiedTime

	return nil
}
