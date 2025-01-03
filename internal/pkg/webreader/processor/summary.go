package processor

import (
	"context"
	"fmt"
	"recally/internal/pkg/config"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader"
	"strings"
	"text/template"
)

const defaultSummaryPrompt = `You are an experienced editor at* **The Wall Street Journal**. *Your task is to read the following article and provide a comprehensive summary for a busy reader who wants to quickly grasp the essential information.

<ResponseFormat>
# Category
[Identify the main category or categories the article belongs to (e.g., Finance, Technology, Health, International News).]

# Summary
[(2-3 sentences) Write a brief abstract summarizing the essence of the article.]

# Abstract
[(approximately 150-200 words) Provide a detailed yet concise summary covering all key information, arguments, and narratives presented in the article.]

# Key Points
[List the most critical points or takeaways in bullet form.]

# Insights and Implications
[Discuss significant insights, implications, or conclusions drawn from the article. Explain how the article relates to broader industry trends or current events.]

# Actionable Takeaways
[(if applicable) Provide any practical advice or recommendations mentioned in the article.]

# Critical Analysis
[Mention any potential biases, assumptions, strengths, or weaknesses in the article. Note any limitations or areas that would benefit from further exploration.]
</ResponseFormat>

Please ensure that your summary is written in the professional, clear, and engaging style characteristic of* **The Wall Street Journal**. *Maintain a neutral and informative tone suitable for helping the reader understand the article without reading it in full.
`

const summaryPromptTemplate = `
{{ .Prompt }}

<Article>
{{ .Article }}
</Article>

<OutputLanguage>
{{or .Language "[Same as article language]"}}
</OutputLanguage>

Please provide your summary below:
`

var summaryPromptTempl = template.Must(template.New("summaryPromptTemplate").Parse(summaryPromptTemplate))

// SummaryOption represents an option for configuring the SummaryProcessor
type SummaryOption func(*SummaryProcessor)

// WithSummaryOptionModel sets the model for the SummaryProcessor
func WithSummaryOptionModel(model string) SummaryOption {
	return func(p *SummaryProcessor) {
		p.config.Model = model
	}
}

// WithSummaryOptionPrompt sets the prompt template for the SummaryProcessor
func WithSummaryOptionPrompt(prompt string) SummaryOption {
	return func(p *SummaryProcessor) {
		p.config.Prompt = prompt
	}
}

func WithSummaryOptionLanguage(language string) SummaryOption {
	return func(p *SummaryProcessor) {
		p.config.Language = language
	}
}

// SummaryConfig contains configuration for the summary processor
type SummaryConfig struct {
	Model    string `json:"model"`
	Prompt   string `json:"prompt"`
	Language string `json:"language"`
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
			Model:  config.Settings.OpenAI.Model, // default model
			Prompt: defaultSummaryPrompt,         // default prompt
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
	var prompt strings.Builder
	if err := summaryPromptTempl.Execute(&prompt, map[string]interface{}{
		"Prompt":   p.config.Prompt,
		"Article":  content.Markwdown,
		"Language": p.config.Language,
	}); err != nil {
		return fmt.Errorf("generate summary prompt: %w", err)
	}

	logger.FromContext(ctx).Info("start summary article", "model", p.config.Model, "language", p.config.Language)
	summary, err := p.llm.TextCompletion(ctx, prompt.String(), llms.WithModel(p.config.Model))
	if err != nil {
		return fmt.Errorf("generate summary: %w", err)
	}

	// Set summary content
	content.Summary = summary

	return nil
}
