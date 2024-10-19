package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/thalesfsp/httpclient/v2"
	"github.com/thalesfsp/inference/internal/config"
	"github.com/thalesfsp/inference/message"
	"github.com/thalesfsp/inference/provider"
	"github.com/thalesfsp/status"
	"github.com/thalesfsp/sypl"
	"github.com/thalesfsp/sypl/level"
	"github.com/thalesfsp/validation"
)

//////
// Const, vars, and types.
//////

// Name of the provider.
const Name = "anthropic"

// Singleton.
var singleton provider.IProvider

// Anthropic provider definition.
type Anthropic struct {
	*provider.Provider

	// Endpoint of the LLM provider.
	Endpoint string `json:"url" validate:"required"`

	// Token of the LLM provider.
	Token string `json:"-" validate:"required"`

	client *httpclient.Client
}

//////
// Implements the IProvider interface.
//////

// Completion generates a completion using the provider API.
// Optionally pass WithResponseBody to unmarshal the response body.
// It will always return the original, unparsed response body, if no error.
//
// NOTE: Not all options are available for all providers.
func (p *Anthropic) Completion(ctx context.Context, options ...provider.Func) (string, error) {
	//////
	// Options initialization.
	//////

	// Prepend the default model to the options.
	options = append(
		[]provider.Func{
			provider.WithModel(p.DefaultModel),
		},
		options...,
	)

	processedOptions, err := provider.NewOptionsFrom(options...)
	if err != nil {
		return "", err
	}

	//////
	// Messages processing.
	//////

	finalMessages := message.NewMessages(
		[]string{},
		processedOptions.UserMessages,
	)

	//////
	// Request body formation.
	//////

	reqBody := &RequestBody{
		Messages: finalMessages,
		Model:    processedOptions.Model,
		Stream:   processedOptions.Stream,

		MaxTokens:   processedOptions.MaxTokens,
		Temperature: processedOptions.Temperature,
		TopP:        processedOptions.TopP,
		TopK:        processedOptions.TopK,
	}

	if len(processedOptions.SystemMessages) > 0 {
		reqBody.System = processedOptions.SystemMessages[0]
	}

	//////
	// Call LLM provider.
	//////

	var respBody ResponseBody

	// Track performance.
	now := time.Now()

	resp, err := p.client.Post(
		ctx,
		p.Endpoint,
		httpclient.WithHeader("x-api-key", p.Token),
		httpclient.WithHeader("anthropic-version", "2023-06-01"),
		httpclient.WithReqBody(reqBody),
		httpclient.WithRespBody(&respBody),
	)
	if err != nil {
		p.GetCounterCompletionFailed().Add(1)

		return "", err
	}

	defer resp.Body.Close()

	// Response processing.
	response, err := ProcessResponse(respBody)
	if err != nil {
		p.GetCounterCompletionFailed().Add(1)

		return "", err
	}

	// Optional response body processing.
	if processedOptions.ResponseBody != nil {
		if err := json.Unmarshal(
			[]byte(response),
			processedOptions.ResponseBody,
		); err != nil {
			return "", err
		}
	}

	//////
	// Observability.
	//////

	// Logging.
	p.GetLogger().PrintlnWithOptions(
		level.Debug,
		fmt.Sprintf("Completion %s", status.Created.String()),
		sypl.WithField("duration", time.Since(now)),
	)

	// Metrics.
	p.GetCounterCompletion().Add(1)

	return response, nil
}

// GetClient returns the client.
func (p *Anthropic) GetClient() any {
	return p.client
}

//////
// Factory.
//////

// New creates a new Anthropic provider.
func New(
	options ...provider.ClientFunc,
) (*Anthropic, error) {
	// Enforces IProvider interface implementation.
	var _ provider.IProvider = (*Anthropic)(nil)

	p, err := provider.New(Name, options...)
	if err != nil {
		return nil, err
	}

	client, err := httpclient.NewDefault(
		httpclient.WithClientName(Name),
	)
	if err != nil {
		return nil, err
	}

	provider := &Anthropic{
		Provider: p,

		Endpoint: p.Endpoint,
		Token:    p.Token,

		client: client,
	}

	if err := validation.Validate(provider); err != nil {
		return nil, err
	}

	singleton = provider

	return provider, nil
}

// NewDefault creates a new Anthropic provider with default values.
func NewDefault(options ...provider.ClientFunc) (*Anthropic, error) {
	opts := []provider.ClientFunc{
		provider.WithEndpoint(config.Get().AnthropicEndpoint),
		provider.WithToken(config.Get().AnthropicToken),
	}

	opts = append(opts, options...)

	return New(opts...)
}

//////
// Exported functionalities.
//////

// Get returns a setup MongoDB, or set it up.
func Get() provider.IProvider {
	if singleton == nil {
		panic(fmt.Sprintf("%s %s not %s", Name, provider.Type, status.Initialized))
	}

	return singleton
}

// Set sets the provider, primarily used for testing.
func Set(s provider.IProvider) {
	singleton = s
}
