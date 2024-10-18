package anthropic

import "github.com/thalesfsp/inference/message"

//////
// Const, vars, types.
//////

//////
// Request body.

// RequestBody represents the request body for the API.
type RequestBody struct {
	Messages []message.Message `json:"messages"`
	Model    string            `json:"model"`
	Stream   bool              `json:"stream"`

	MaxTokens   int     `json:"max_tokens,omitempty"`
	System      string  `json:"system,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	TopK        int     `json:"top_k,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
}

//////
// Response body.

// Usage definition.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// Content definition.
type Content struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

// ResponseBody represents the response body from the API.
type ResponseBody struct {
	Content      []Content `json:"content"`
	ID           string    `json:"id"`
	Model        string    `json:"model"`
	Role         string    `json:"role"`
	StopReason   string    `json:"stop_reason"`
	StopSequence any       `json:"stop_sequence"`
	Type         string    `json:"type"`
	Usage        Usage     `json:"usage"`
}
