package webreader

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"recally/internal/pkg/logger"
	"time"

	"github.com/go-shiori/go-readability"
)

// Content represents the result of a fetch operation
type Content struct {
	URL         string `json:"url"`
	Content     io.ReadCloser
	StatusCode  int
	ContentType string
	Headers     map[string][]string

	Title         string     `json:"title"`
	Html          string     `json:"html"`
	Text          string     `json:"text"`
	Markwdown     string     `json:"markdown"`
	Summary       string     `json:"summary"`
	Image         string     `json:"image"`
	Description   string     `json:"description"`
	PublishedTime *time.Time `json:"published_time"`
	ModifiedTime  *time.Time `json:"modified_time"`
}

// Fetcher defines the interface for different content fetchers
type Fetcher interface {
	// Fetch retrieves content from the given URL
	Fetch(ctx context.Context, url string) (*Content, error)
	// Close cleans up any resources used by the fetcher
	Close() error
}

// Processor defines the interface for different content processors
type Processor interface {
	// Process processes the input content string and returns the processed result string
	Process(ctx context.Context, input *Content) error

	Name() string
}

// Reader represents a configurable web content reader
type Reader struct {
	fetcher    Fetcher
	processors []Processor
}

// New creates a new Reader instance
func New(f Fetcher, processors ...Processor) *Reader {
	return &Reader{
		fetcher:    f,
		processors: processors,
	}
}

// AddProcessor adds a new processor to the end of the processing chain
func (w *Reader) AddProcessor(p Processor) {
	w.processors = append(w.processors, p)
}

// Read fetches content from the URL and processes it through all processors in sequence
func (w *Reader) Read(ctx context.Context, url string) (*Content, error) {
	// Fetch the content
	defer w.fetcher.Close()
	content, err := w.fetcher.Fetch(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("fetch error: %w", err)
	}

	// Pre-process the content
	if err := w.preProcess(content); err != nil {
		return nil, fmt.Errorf("pre-process error: %w", err)
	}

	// Post Process the content through each processor in sequence
	for _, p := range w.processors {
		if err := p.Process(ctx, content); err != nil {
			logger.FromContext(ctx).Error("process error at processor", "processor", p.Name(), "err", err)
		}
	}

	return content, nil
}

func (w *Reader) preProcess(content *Content) error {
	parsedURL, err := url.ParseRequestURI(content.URL)
	// If there's an error parsing the URI, set parsedURL to nil
	if err != nil {
		parsedURL = nil
	}
	defer content.Content.Close()
	// Use the readability package's FromReader function to parse the HTML content
	article, err := readability.FromReader(content.Content, parsedURL)
	// If there's an error parsing the HTML content, return the error
	if err != nil {
		return fmt.Errorf("failed to parse %s, %v", content.URL, err)
	}

	// Set Markdown content
	content.Title = article.Title
	content.Description = article.Byline
	if article.Favicon != "" {
		content.Image = article.Favicon
	}
	if content.Image == "" {
		content.Image = article.Image
	}
	content.Text = article.TextContent
	content.Html = article.Content
	content.PublishedTime = article.PublishedTime
	content.ModifiedTime = article.ModifiedTime
	return nil
}
