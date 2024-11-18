package processor

import (
	"context"
	"fmt"
	"strings"
	"vibrain/internal/pkg/llms"
	"vibrain/internal/pkg/webreader"
)

const summaryPrompt = `
You are an expert analyst skilled in critical reading and synthesis. Approach this text with careful analytical thinking to create an insightful, well-structured summary.

Instructions:
1. Read the text thoroughly, analyzing:
- The fundamental argument or thesis
- The logical structure of key supporting points
- The quality and application of evidence
- The author's methodology and reasoning
- Significant implications and conclusions
- Any underlying assumptions or limitations

2. Synthesize your analysis into this exact format:

# Summary
[Provide a thorough 3-5 paragraph analytical summary that:
- Opens with a clear statement of the main argument/thesis
- Explains how key points build and connect to support the main argument
- Evaluates the strength and relevance of evidence
- Examines the author's reasoning and methodology
- Discusses significant implications and conclusions
Each paragraph should flow logically into the next, creating a coherent analysis rather than just a list of points.]

# Opinions
[Leave blank - retain heading for format]

<inputText>
{inputText}
</inputText>
`

// SummaryOption represents an option for configuring the SummaryProcessor
type SummaryOption func(*SummaryProcessor)

// WithModel sets the model for the SummaryProcessor
func WithModel(model string) SummaryOption {
	return func(p *SummaryProcessor) {
		p.config.Model = model
	}
}

// WithPrompt sets the prompt template for the SummaryProcessor
func WithPrompt(prompt string) SummaryOption {
	return func(p *SummaryProcessor) {
		p.config.Prompt = prompt
	}
}

// SummaryConfig contains configuration for the summary processor
type SummaryConfig struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// SummaryProcessor implements content summarization using LLM
type SummaryProcessor struct {
	config SummaryConfig
	llm    *llms.LLM
}

// NewSummaryProcessor creates a new SummaryProcessor with the given configuration
func NewSummaryProcessor(llm *llms.LLM, opts ...SummaryOption) *SummaryProcessor {
	p := &SummaryProcessor{
		config: SummaryConfig{
			Model:  "gpt-4-mini",  // default model
			Prompt: summaryPrompt, // default prompt
		},
		llm: llm,
	}

	// Apply options
	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *SummaryProcessor) Name() string {
	return "AI Summary"
}

// Process implements the Processor interface
func (p *SummaryProcessor) Process(ctx context.Context, content *webreader.Content) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	prompt := strings.ReplaceAll(p.config.Prompt, "{inputText}", content.Html)

	summary, err := p.llm.TextCompletion(ctx, prompt, llms.WithModel(p.config.Model))
	if err != nil {
		return fmt.Errorf("generate summary: %w", err)
	}

	// Set summary content
	content.Summary = summary

	return nil
}
