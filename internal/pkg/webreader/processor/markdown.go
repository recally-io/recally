package processor

import (
	"context"
	"fmt"
	"recally/internal/pkg/webreader"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

// MarkdownOption represents an option for configuring the MarkdownProcessor
type MarkdownOption func(*MarkdownProcessor)

// WithMarkdownOptions sets the markdown converter options
func WithMarkdownOptions(options md.Options) MarkdownOption {
	return func(p *MarkdownProcessor) {
		p.config.Options = options
	}
}

// MarkdownConfig contains configuration for the markdown processor
type MarkdownConfig struct {
	Options md.Options
}

// MarkdownProcessor implements HTML to Markdown conversion
type MarkdownProcessor struct {
	config MarkdownConfig
	conv   *md.Converter
}

// NewMarkdownProcessor creates a new MarkdownProcessor with the given configuration
func NewMarkdownProcessor(opts ...MarkdownOption) *MarkdownProcessor {
	p := &MarkdownProcessor{
		config: MarkdownConfig{
			Options: md.Options{}, // default options
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(p)
	}

	// Create converter with configured options
	p.conv = md.NewConverter("", true, &p.config.Options)

	return p
}

func (p *MarkdownProcessor) Name() string {
	return "Markdown"
}

// Process implements the Processor interface
func (p *MarkdownProcessor) Process(ctx context.Context, content *webreader.Content) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Convert HTML to Markdown
	markdown, err := p.conv.ConvertString(content.Html)
	if err != nil {
		return fmt.Errorf("webreader markdown converter error: %w", err)
	}

	// Set Markdown content
	content.Markwdown = markdown
	return nil
}
