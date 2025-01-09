package webreader

import (
	"context"
	"fmt"
	"io"
	"recally/internal/pkg/logger"
	"time"
)

// Content represents the result of a fetch operation
type Content struct {
	URL           string     `json:"url"`
	Title         string     `json:"title"`
	Html          string     `json:"html"`
	Text          string     `json:"text"`
	Markwdown     string     `json:"markdown"`
	Summary       string     `json:"summary"`
	Cover         string     `json:"cover"`
	Favicon       string     `json:"favicon"`
	Image         string     `json:"image"`
	Author        string     `json:"author"`
	Description   string     `json:"description"`
	SiteName      string     `json:"site_name"`
	PublishedTime *time.Time `json:"published_time"`
	ModifiedTime  *time.Time `json:"modified_time"`
}

type FetchedContent struct {
	Content

	Reader      io.ReadCloser // Content is the raw content of the fetched URL
	StatusCode  int
	ContentType string
	Headers     map[string][]string
}

// Fetcher defines the interface for different content fetchers
type Fetcher interface {
	// Fetch retrieves content from the given URL
	Fetch(ctx context.Context, url string) (*FetchedContent, error)
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
	fetchedContent, err := w.fetcher.Fetch(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("fetch error: %w", err)
	}

	// Pre-process the content
	if err := w.preProcess(fetchedContent); err != nil {
		return nil, fmt.Errorf("pre-process error: %w", err)
	}

	content := &fetchedContent.Content
	// Post Process the content through each processor in sequence
	for _, p := range w.processors {
		if err := p.Process(ctx, content); err != nil {
			logger.FromContext(ctx).Error("process error at processor", "processor", p.Name(), "err", err)
		}
	}

	return content, nil
}

func (w *Reader) preProcess(content *FetchedContent) error {
	rawContent, err := io.ReadAll(content.Reader)
	if err != nil {
		return fmt.Errorf("failed to read content: %w", err)
	}
	defer content.Reader.Close()

	if content.Html == "" {
		content.Html = string(rawContent)
	}

	return nil
}
