package processor

import (
	"context"
	"fmt"
	"recally/internal/pkg/webreader"
	"recally/internal/pkg/webreader/processor/hooks"

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

func WithMarkdownBeforeHook(hooks ...md.BeforeHook) MarkdownOption {
	return func(p *MarkdownProcessor) {
		p.config.BeforeHooks = append(p.config.BeforeHooks, hooks...)
	}
}

func WithMarkdownAfterHook(hooks ...md.Afterhook) MarkdownOption {
	return func(p *MarkdownProcessor) {
		p.config.AfterHooks = append(p.config.AfterHooks, hooks...)
	}
}

// MarkdownConfig contains configuration for the markdown processor
type MarkdownConfig struct {
	Options     md.Options
	BeforeHooks []md.BeforeHook
	AfterHooks  []md.Afterhook
}

// MarkdownProcessor implements HTML to Markdown conversion
type MarkdownProcessor struct {
	Host   string `json:"host"`
	config MarkdownConfig
	conv   *md.Converter
}

// NewMarkdownProcessor creates a new MarkdownProcessor with the given configuration
func NewMarkdownProcessor(host string, opts ...MarkdownOption) *MarkdownProcessor {
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
	p.conv = md.NewConverter(host, true, &p.config.Options)

	// Register buildin hooks
	p.conv.Before(hooks.GetMarkdownBeforeHooks(host)...)
	p.conv.After(hooks.GetMarkdownAfterHooks(host)...)

	// Register custom hooks
	if len(p.config.BeforeHooks) > 0 {
		p.conv.Before(p.config.BeforeHooks...)
	}
	if len(p.config.AfterHooks) > 0 {
		p.conv.After(p.config.AfterHooks...)
	}

	return p
}

func (p *MarkdownProcessor) Name() string {
	return "Markdown"
}

// Process implements the Processor interface
func (p *MarkdownProcessor) Process(ctx context.Context, content *webreader.Content) error {
	// Convert HTML to Markdown
	markdown, err := p.conv.ConvertString(content.Html)
	if err != nil {
		return fmt.Errorf("webreader markdown converter error: %w", err)
	}

	// Set Markdown content
	content.Markwdown = markdown
	return nil
}
