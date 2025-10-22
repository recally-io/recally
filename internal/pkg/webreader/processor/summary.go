package processor

import (
	"context"
	"fmt"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/config"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader"
	"strings"
	"text/template"
)

const defaultSummaryPrompt = `You are an experienced editor at **The Wall Street Journal**. Your task is to read the following article and provide a comprehensive summary for a busy reader who wants to quickly grasp the essential information.

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

Please ensure that your summary is written in the professional, clear, and engaging style characteristic of **The Wall Street Journal**. Maintain a neutral and informative tone suitable for helping the reader understand the article without reading it in full.
`

const summaryPromptTemplate = `
<Instruction>
{{ .Prompt }}
</Instruction>

<Article>
{{ .Article }}
</Article>

Please provide your summary using {{or .Language "[Same as article language]"}}  below. When you finish summarizing the article, please add tags to your response in following format:

<OutputFormat>
<summary>
[Your comprehensive summary that follow the Instruction]
</summary>

<tags>
[Comma-Separated Tags here]
</tags>
</OutputFormat>
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

func WithSummaryOptionUser(user *auth.UserDTO) SummaryOption {
	return func(p *SummaryProcessor) {
		if user.Settings.SummaryOptions.Model != "" {
			p.config.Model = user.Settings.SummaryOptions.Model
		}
		if user.Settings.SummaryOptions.Prompt != "" {
			p.config.Prompt = user.Settings.SummaryOptions.Prompt
		}
		if user.Settings.SummaryOptions.Language != "" {
			p.config.Language = user.Settings.SummaryOptions.Language
		}
	}
}

// SummaryProcessor implements content summarization using LLM
type SummaryProcessor struct {
	config auth.AIConfig
	llm    *llms.LLM
}

// NewSummaryProcessor creates a new SummaryProcessor with the given configuration
func NewSummaryProcessor(llm *llms.LLM, opts ...SummaryOption) *SummaryProcessor {
	p := &SummaryProcessor{
		config: auth.AIConfig{
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
	prompt, err := p.buildPrompt(ctx, content.Markwdown)
	if err != nil {
		return err
	}
	logger.FromContext(ctx).Info("start summary article", "model", p.config.Model, "language", p.config.Language, "streaming", false)
	summary, err := p.llm.TextCompletion(ctx, prompt, llms.WithModel(p.config.Model))
	if err != nil {
		return fmt.Errorf("generate summary: %w", err)
	}

	// Set summary content
	content.Summary = summary

	return nil
}

func (p *SummaryProcessor) StreamingSummary(ctx context.Context, content string, streamingFunc func(content llms.StreamingString)) {
	prompt, err := p.buildPrompt(ctx, content)
	if err != nil {
		streamingFunc(llms.StreamingString{
			Err: err,
		})
		return
	}
	logger.FromContext(ctx).Info("start streaming summary article", "model", p.config.Model, "language", p.config.Language, "streaming", true)
	p.llm.StreamingTextCompletion(ctx, prompt, streamingFunc, llms.WithModel(p.config.Model))
}

func (p *SummaryProcessor) buildPrompt(ctx context.Context, content string) (string, error) {
	var prompt strings.Builder
	if err := summaryPromptTempl.Execute(&prompt, map[string]any{
		"Prompt":   p.config.Prompt,
		"Article":  content,
		"Language": p.config.Language,
	}); err != nil {
		logger.FromContext(ctx).Error("error generate summary prompt", "err", err)
		return "", fmt.Errorf("generate summary prompt: %w", err)
	}
	return prompt.String(), nil
}

func (p *SummaryProcessor) ParseSummaryInfo(content string) (summary string, tags []string) {
	summary = parseXmlContent(content, "summary")
	tagString := parseXmlContent(content, "tags")
	tags = tagStringToArray(tagString)
	return
}
