package embedddings

import "github.com/sashabaranov/go-openai"

const defaultModel = "text-embedding-3-small"

// Option is a function that configures a CallOptions.
type Option func(*Options)

// Options is a set of options for calling models. Not all models support
// all options.
type Options struct {
	// Model is the model to use.
	Model string
}

func (o Options) ToEmbeddingRequest() openai.EmbeddingRequestStrings {
	if o.Model == "" {
		o.Model = defaultModel
	}

	return openai.EmbeddingRequestStrings{
		Model: openai.EmbeddingModel(o.Model),
	}
}

// WithModel specifies which model name to use.
func WithModel(model string) Option {
	return func(o *Options) {
		o.Model = model
	}
}
