package processor

import (
	"context"
	"fmt"
	"net/url"
	"recally/internal/pkg/webreader"
	"recally/internal/pkg/webreader/processor/hooks"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

	// Use the readability package's FromReader function to parse the HTML content
	article, err := readability.FromReader(strings.NewReader(content.Html), parsedURL)
	// If there's an error parsing the HTML content, return the error
	if err != nil {
		return fmt.Errorf("failed to parse %s, %v", content.URL, err)
	}

	// Set meta info
	content.Title = article.Title
	content.Author = article.Byline
	content.Description = article.Excerpt
	content.SiteName = article.SiteName

	// Set the cover image
	content.Cover = article.Image
	content.Favicon = article.Favicon

	// set content
	content.Text = article.TextContent

	// Set the published and modified time
	content.PublishedTime = article.PublishedTime
	content.ModifiedTime = article.ModifiedTime

	// run hooks
	hooks := hooks.GetReadabilityHooks(parsedURL.Host)
	if len(hooks) > 0 {
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(content.Html))
		for _, hook := range hooks {
			hook(doc, content)
		}
	}

	return nil
}
