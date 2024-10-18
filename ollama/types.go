package ollama

import (
	"time"

	"github.com/thalesfsp/inference/message"
)

//////
// Const, vars, types.
//////

//////
// Request body.

// RequestBodyOptions represents the options for the request body.
type RequestBodyOptions struct {
	Seed        int     `json:"seed,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	TopK        int     `json:"top_k,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
}

// RequestBody represents the request body for the Ollama API.
type RequestBody struct {
	Messages []message.Message `json:"messages"`
	Model    string            `json:"model"`
	Stream   bool              `json:"stream"`

	Options RequestBodyOptions `json:"options,omitempty"`
}

//////
// Response body.

// ResponseBody represents the response body from the Ollama API.
type ResponseBody struct {
	Context            []int           `json:"context,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	Done               bool            `json:"done"`
	DoneReason         string          `json:"done_reason"`
	EvalCount          int             `json:"eval_count"`
	EvalDuration       int64           `json:"eval_duration"`
	LoadDuration       int             `json:"load_duration"`
	Message            message.Message `json:"message"`
	Model              string          `json:"model"`
	PromptEvalCount    int             `json:"prompt_eval_count"`
	PromptEvalDuration int             `json:"prompt_eval_duration"`
	TotalDuration      int64           `json:"total_duration"`
}
