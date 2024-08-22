package llms

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

const (
	defaultModel     = OpenAIGPT4oMini
	ToolTypeFunction = "function"
)

// Option is a function that configures a CallOptions.
type Option func(*Options)

// Options is a set of options for calling models. Not all models support
// all options.
type Options struct {
	// Model is the model to use.
	Model string `json:"model"`
	// CandidateCount is the number of response candidates to generate.
	CandidateCount int `json:"candidate_count"`
	// MaxTokens is the maximum number of tokens to generate.
	MaxTokens int `json:"max_tokens"`
	// Temperature is the temperature for sampling, between 0 and 1.
	Temperature float64 `json:"temperature"`
	// StopWords is a list of words to stop on.
	StopWords []string `json:"stop_words"`
	// StreamingFunc is a function to be called for each chunk of a streaming response.
	// Return an error to stop streaming early.
	StreamingFunc func(ctx context.Context, chunk []byte) error `json:"-"`
	// TopK is the number of tokens to consider for top-k sampling.
	TopK int `json:"top_k"`
	// TopP is the cumulative probability for top-p sampling.
	TopP float64 `json:"top_p"`
	// Seed is a seed for deterministic sampling.
	Seed int `json:"seed"`
	// MinLength is the minimum length of the generated text.
	MinLength int `json:"min_length"`
	// MaxLength is the maximum length of the generated text.
	MaxLength int `json:"max_length"`
	// N is how many chat completion choices to generate for each input message.
	N int `json:"n"`
	// RepetitionPenalty is the repetition penalty for sampling.
	RepetitionPenalty float64 `json:"repetition_penalty"`
	// FrequencyPenalty is the frequency penalty for sampling.
	FrequencyPenalty float64 `json:"frequency_penalty"`
	// PresencePenalty is the presence penalty for sampling.
	PresencePenalty float64 `json:"presence_penalty"`

	// JSONMode is a flag to enable JSON mode.
	JSONMode bool `json:"json"`

	// Tools is a list of tools to use. Each tool can be a specific tool or a function.
	Tools     []Tool   `json:"tools,omitempty"`
	ToolNames []string `json:"tool_names,omitempty"`
	// ToolChoice is the choice of tool to use, it can either be "none", "auto" (the default behavior), or a specific tool as described in the ToolChoice type.
	ToolChoice any `json:"tool_choice"`

	// Function defitions to include in the request.
	// Deprecated: Use Tools instead.
	Functions []FunctionDefinition `json:"functions,omitempty"`
	// FunctionCallBehavior is the behavior to use when calling functions.
	//
	// If a specific function should be invoked, use the format:
	// `{"name": "my_function"}`
	// Deprecated: Use ToolChoice instead.
	FunctionCallBehavior FunctionCallBehavior `json:"function_call,omitempty"`

	// Metadata is a map of metadata to include in the request.
	// The meaning of this field is specific to the backend in use.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func (o Options) ToChatCompletionRequest() openai.ChatCompletionRequest {
	if o.Model == "" {
		o.Model = defaultModel
	}
	return openai.ChatCompletionRequest{
		Model:            o.Model,
		Stop:             o.StopWords,
		Temperature:      float32(o.Temperature),
		MaxTokens:        o.MaxTokens,
		N:                o.N,
		FrequencyPenalty: float32(o.FrequencyPenalty),
		PresencePenalty:  float32(o.PresencePenalty),

		ToolChoice: o.ToolChoice,
		Seed:       &o.Seed,
	}
}

// Tool is a tool that can be used by the model.
type Tool struct {
	// Type is the type of the tool.
	Type string `json:"type"`
	// Function is the function to call.
	Function *FunctionDefinition `json:"function,omitempty"`
}

// FunctionDefinition is a definition of a function that can be called by the model.
type FunctionDefinition struct {
	// Name is the name of the function.
	Name string `json:"name"`
	// Description is a description of the function.
	Description string `json:"description"`
	// Parameters is a list of parameters for the function.
	Parameters any `json:"parameters,omitempty"`
}

// ToolChoice is a specific tool to use.
type ToolChoice struct {
	// Type is the type of the tool.
	Type string `json:"type"`
	// Function is the function to call (if the tool is a function).
	Function *FunctionReference `json:"function,omitempty"`
}

// FunctionReference is a reference to a function.
type FunctionReference struct {
	// Name is the name of the function.
	Name string `json:"name"`
}

// FunctionCallBehavior is the behavior to use when calling functions.
type FunctionCallBehavior string

const (
	// FunctionCallBehaviorNone will not call any functions.
	FunctionCallBehaviorNone FunctionCallBehavior = "none"
	// FunctionCallBehaviorAuto will call functions automatically.
	FunctionCallBehaviorAuto FunctionCallBehavior = "auto"
)

// WithModel specifies which model name to use.
func WithModel(model string) Option {
	return func(o *Options) {
		o.Model = model
	}
}

// WithMaxTokens specifies the max number of tokens to generate.
func WithMaxTokens(maxTokens int) Option {
	return func(o *Options) {
		o.MaxTokens = maxTokens
	}
}

// WithCandidateCount specifies the number of response candidates to generate.
func WithCandidateCount(c int) Option {
	return func(o *Options) {
		o.CandidateCount = c
	}
}

// WithTemperature specifies the model temperature, a hyperparameter that
// regulates the randomness, or creativity, of the AI's responses.
func WithTemperature(temperature float64) Option {
	return func(o *Options) {
		o.Temperature = temperature
	}
}

// WithStopWords specifies a list of words to stop generation on.
func WithStopWords(stopWords []string) Option {
	return func(o *Options) {
		o.StopWords = stopWords
	}
}

// WithOptions specifies options.
func WithOptions(options Options) Option {
	return func(o *Options) {
		(*o) = options
	}
}

// WithStreamingFunc specifies the streaming function to use.
func WithStreamingFunc(streamingFunc func(ctx context.Context, chunk []byte) error) Option {
	return func(o *Options) {
		o.StreamingFunc = streamingFunc
	}
}

// WithTopK will add an option to use top-k sampling.
func WithTopK(topK int) Option {
	return func(o *Options) {
		o.TopK = topK
	}
}

// WithTopP	will add an option to use top-p sampling.
func WithTopP(topP float64) Option {
	return func(o *Options) {
		o.TopP = topP
	}
}

// WithSeed will add an option to use deterministic sampling.
func WithSeed(seed int) Option {
	return func(o *Options) {
		o.Seed = seed
	}
}

// WithMinLength will add an option to set the minimum length of the generated text.
func WithMinLength(minLength int) Option {
	return func(o *Options) {
		o.MinLength = minLength
	}
}

// WithMaxLength will add an option to set the maximum length of the generated text.
func WithMaxLength(maxLength int) Option {
	return func(o *Options) {
		o.MaxLength = maxLength
	}
}

// WithN will add an option to set how many chat completion choices to generate for each input message.
func WithN(n int) Option {
	return func(o *Options) {
		o.N = n
	}
}

// WithRepetitionPenalty will add an option to set the repetition penalty for sampling.
func WithRepetitionPenalty(repetitionPenalty float64) Option {
	return func(o *Options) {
		o.RepetitionPenalty = repetitionPenalty
	}
}

// WithFrequencyPenalty will add an option to set the frequency penalty for sampling.
func WithFrequencyPenalty(frequencyPenalty float64) Option {
	return func(o *Options) {
		o.FrequencyPenalty = frequencyPenalty
	}
}

// WithPresencePenalty will add an option to set the presence penalty for sampling.
func WithPresencePenalty(presencePenalty float64) Option {
	return func(o *Options) {
		o.PresencePenalty = presencePenalty
	}
}

// WithTools will add an option to set the tools to use.
func WithToolNames(names []string) Option {
	return func(o *Options) {
		o.ToolNames = names
	}
}

// WithJSONMode will add an option to set the response format to JSON.
// This is useful for models that return structured data.
func WithJSONMode() Option {
	return func(o *Options) {
		o.JSONMode = true
	}
}
