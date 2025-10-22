package processor

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"recally/internal/core/files"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/config"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader"

	"github.com/sashabaranov/go-openai"
)

const defaultSummaryImagePrompt = `You are an expert image analyst with strong skills in visual interpretation and metadata generation. When provided with an image, generate:

1. **Title**: Create a concise, engaging title (3-8 words) that captures the image's essence.
2. **Description**: Write a detailed 2-4 sentence description covering:
   - Key visual elements (objects, people, scenery)
   - Colors, lighting, and artistic style
   - Atmosphere/mood
   - Notable details or focal points
3. **Tags**: List 3-5 relevant keywords or phrases (separated by commas) including:
   - Main subjects
   - Colors/palette
   - Style (e.g., photorealistic, abstract)
   - Themes/concepts

<Guidelines>
- Focus only on observable elements (avoid assumptions about context)
- Prioritize clarity and accuracy over creativity
- Use neutral, objective language
</Guidelines>

<OutputFormat>
<title>
[Title here]
</title>

<description>
[Description here]
</description>

<tags>
[Comma-Separated Tags here]
</tags>
</OutputFormat>

<ExampleOutput>
<title>
Sunset Over Mountain Lake
</title>

<description>
A serene alpine lake reflects vibrant orange and pink sunset hues, surrounded by pine-covered slopes. The hyper-realistic digital painting features crisp water reflections and dramatic cloud formations, creating a peaceful yet awe-inspiring atmosphere.
</description>

<tags>
landscape, sunset, lake, mountains, digital painting
</tags>
</ExampleOutput>
`

// SummaryOption represents an option for configuring the SummaryProcessor
type SummaryImageOption func(*SummaryImageProcessor)

// WithSummaryOptionModel sets the model for the SummaryProcessor
func WithSummaryImageOptionModel(model string) SummaryImageOption {
	return func(p *SummaryImageProcessor) {
		p.config.Model = model
	}
}

// WithSummaryOptionPrompt sets the prompt template for the SummaryProcessor
func WithSummaryImageOptionPrompt(prompt string) SummaryImageOption {
	return func(p *SummaryImageProcessor) {
		p.config.Prompt = prompt
	}
}

func WithSummaryImageOptionLanguage(language string) SummaryImageOption {
	return func(p *SummaryImageProcessor) {
		p.config.Language = language
	}
}

func WithSummaryImageOptionUser(user *auth.UserDTO) SummaryImageOption {
	return func(p *SummaryImageProcessor) {
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
type SummaryImageProcessor struct {
	config auth.AIConfig
	llm    *llms.LLM
}

// NewSummaryProcessor creates a new SummaryProcessor with the given configuration
func NewSummaryImageProcessor(llm *llms.LLM, opts ...SummaryImageOption) *SummaryImageProcessor {
	p := &SummaryImageProcessor{
		config: auth.AIConfig{
			Model:  config.Settings.OpenAI.VisionModel, // default model
			Prompt: defaultSummaryImagePrompt,          // default prompt
		},
		llm: llm,
	}

	// Apply options
	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *SummaryImageProcessor) Name() string {
	return "AI Image Summarization"
}

func (p *SummaryImageProcessor) Process(ctx context.Context, content *webreader.Content) error {
	return nil
}

func (p *SummaryImageProcessor) StreamingSummary(ctx context.Context, imgURL string, streamingFunc func(content llms.StreamingMessage)) {
	p.process(ctx, imgURL, streamingFunc, true)
}

func (p *SummaryImageProcessor) Summary(ctx context.Context, imgURL string, streamingFunc func(content llms.StreamingMessage)) {
	p.process(ctx, imgURL, streamingFunc, false)
}

// Process implements the Processor interface
func (p *SummaryImageProcessor) process(ctx context.Context, imgURL string, streamingFunc func(content llms.StreamingMessage), streaming bool) {
	logger.FromContext(ctx).Info("start describe image", "model", p.config.Model, "language", p.config.Language, "streaming", streaming)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: p.config.Prompt,
		},
		{
			Role: openai.ChatMessageRoleUser,
			MultiContent: []openai.ChatMessagePart{
				{
					Type: openai.ChatMessagePartTypeText,
					Text: fmt.Sprintf("Describe the image in %s", p.config.Language),
				},
				{
					Type: openai.ChatMessagePartTypeImageURL,
					ImageURL: &openai.ChatMessageImageURL{
						URL: imgURL,
					},
				},
			},
		},
	}

	p.llm.GenerateContent(ctx, messages, streamingFunc, llms.WithModel(p.config.Model), llms.WithStream(streaming))
}

func (p *SummaryImageProcessor) EncodeImage(reader io.ReadCloser, fileName string) (string, error) {
	defer func() { _ = reader.Close() }()
	contentType := files.GetFileMIMEWithDefault(fileName, "image/jpeg")
	photoBytes, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read photo: %w", err)
	}
	photoBase64 := base64.StdEncoding.EncodeToString(photoBytes)
	imgUrl := fmt.Sprintf("data:%s;base64,%s", contentType, photoBase64)

	return imgUrl, nil
}

func (p *SummaryImageProcessor) ParseSummaryInfo(content string) (title, description string, tags []string) {
	title = parseXmlContent(content, "title")
	description = parseXmlContent(content, "description")
	tagString := parseXmlContent(content, "tags")
	tags = tagStringToArray(tagString)
	return
}
