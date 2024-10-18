package provider

import (
	"github.com/thalesfsp/validation"
)

//////
// Vars, consts, and types.
//////

// Func allows to set options.
type Func func(o *Options) error

// Options for operations.
type Options struct {
	// Model is the model to be used.
	Model string `json:"model" validate:"required"`

	// UserMessages is the user role messages.
	UserMessages []string `json:"userMessages" validate:"required"`

	// MaxTokens defines the max amount of tokens in the response. Default to 0
	// which means no limit.
	MaxTokens int `json:"maxToken,omitempty" validate:"gte=0"`

	// Seed the LLM will make a best effort to sample deterministically, such
	// that repeated requests with the same seed and parameters should return
	// the same result. Default to not set which means no determinism.
	//
	// NOTE: Determinism is not guaranteed.
	Seed int `json:"seed,omitempty" validate:"gte=0"`

	// Stream determines if the response should be streamed or not. Default to
	// false which means not to stream.
	Stream bool `json:"stream"`

	// SystemMessages is the system role messages.
	SystemMessages []string `json:"systemMessages,omitempty"`

	// Temperature to use, between 0 and 2. Higher values like 0.8 will make the
	// output more random, while lower values like 0.2 will make it more focused
	// and deterministic. Default to 1.0 which means be creative.
	//
	// NOTE: It's generally recommended altering this OR temperature but not
	// both!
	Temperature float64 `json:"temperature,omitempty" validate:"gte=0"`

	// TopK means the number of highest probability vocabulary tokens to keep
	// for sampling, default is not set (0) which means no restrictions.
	TopK int `json:"topK,omitempty" validate:"gte=0"`

	// TopP is an alternative to sampling with temperature, called nucleus
	// sampling, where the model considers the results of the tokens with top_p
	// probability mass. So 0.1 means only the tokens comprising the top 10%
	// probability mass are considered. Default to 1.0 which means consider all
	// tokens.
	//
	// NOTE: It's generally recommended altering this OR temperature but not
	// both!
	TopP float64 `json:"topP,omitempty" validate:"gte=0"`
}

//////
// Exported built-in options.
//////

// WithMaxToken sets the maxToken option.
func WithMaxToken(maxToken int) Func {
	return func(o *Options) error {
		if maxToken > 0 {
			o.MaxTokens = maxToken
		}

		return nil
	}
}

// WithModel sets the model option.
func WithModel(model string) Func {
	return func(o *Options) error {
		if model != "" {
			o.Model = model
		}

		return nil
	}
}

// WithSeed sets the seed option.
func WithSeed(seed int) Func {
	return func(o *Options) error {
		if seed > 0 {
			o.Seed = seed
		}

		return nil
	}
}

// WithTemperature sets the temperature option.
func WithTemperature(temperature float64) Func {
	return func(o *Options) error {
		if temperature > 0 {
			o.Temperature = temperature
		}

		return nil
	}
}

// WithTopK sets the topK option.
func WithTopK(topK int) Func {
	return func(o *Options) error {
		if topK > 0 {
			o.TopK = topK
		}

		return nil
	}
}

// WithTopP sets the topP option.
func WithTopP(topP float64) Func {
	return func(o *Options) error {
		if topP > 0 {
			o.TopP = topP
		}

		return nil
	}
}

// WithStream sets the stream option.
func WithStream(stream bool) Func {
	return func(o *Options) error {
		o.Stream = stream

		return nil
	}
}

// WithSystemMessages sets the systemMessages option.
func WithSystemMessages(systemMessage ...string) Func {
	return func(o *Options) error {
		if len(o.SystemMessages) > 0 {
			// Initialize the slice.
			o.SystemMessages = make([]string, 0)
		}

		o.SystemMessages = append(o.SystemMessages, systemMessage...)

		return nil
	}
}

// WithUserMessages sets the userMessages option.
func WithUserMessages(userMessage ...string) Func {
	return func(o *Options) error {
		if len(o.UserMessages) > 0 {
			// Initialize the slice.
			o.UserMessages = make([]string, 0)
		}

		o.UserMessages = append(o.UserMessages, userMessage...)

		return nil
	}
}

//////
// Factory.
//////

// NewOptionsFrom process, and validate against the default options.
//
//nolint:mnd,gomnd
func NewOptionsFrom(options ...Func) (*Options, error) {
	defaultOptions := Options{
		MaxTokens:   4096,
		Stream:      false,
		Temperature: 0.7,
	}

	// Apply options against the default options.
	for _, option := range options {
		if err := option(&defaultOptions); err != nil {
			return nil, err
		}
	}

	if err := validation.Validate(&defaultOptions); err != nil {
		return nil, err
	}

	return &defaultOptions, nil
}
