package openai

import "github.com/thalesfsp/inference/message"

//////
// Const, vars, types.
//////

//////
// Request body.

// RequestBody represents the request body for the OpenAI API.
type RequestBody struct {
	Messages []message.Message `json:"messages"`
	Model    string            `json:"model"`
	Stream   bool              `json:"stream"`

	MaxTokens   int     `json:"max_completion_tokens,omitempty"`
	Seed        int     `json:"seed,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
}

//////
// Response body.

// Usage OpenAI API definition.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Choice OpenAI API definition.
type Choice struct {
	Message      message.Message `json:"message"`
	FinishReason string          `json:"finish_reason"`
	Index        int             `json:"index"`
}

// ResponseBody represents the response body from the OpenAI API.
type ResponseBody struct {
	Choices []Choice `json:"choices" validate:"gt=0"`
	Created int      `json:"created"`
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Object  string   `json:"object"`
	Usage   Usage    `json:"usage"`
}
