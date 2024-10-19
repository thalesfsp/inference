package ollama

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
const Name = "ollama"

// Singleton.
var singleton provider.IProvider

// Ollama provider definition.
type Ollama struct {
	*provider.Provider

	// Endpoint of the LLM provider.
	Endpoint string `json:"url" validate:"required"`

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
func (p *Ollama) Completion(ctx context.Context, options ...provider.Func) (string, error) {
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
		processedOptions.SystemMessages,
		processedOptions.UserMessages,
	)

	//////
	// Request body formation.
	//////

	reqBody := &RequestBody{
		Messages: finalMessages,
		Model:    processedOptions.Model,
		Stream:   processedOptions.Stream,

		Options: RequestBodyOptions{
			Temperature: processedOptions.Temperature,
			TopK:        processedOptions.TopK,
			TopP:        processedOptions.TopP,
		},
	}

	//////
	// Call LLM provider.
	//////

	var respBody ResponseBody

	// Track performance.
	now := time.Now()

	// Actual call.
	resp, err := p.client.Post(
		ctx,
		p.Endpoint,
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
func (p *Ollama) GetClient() any {
	return p.client
}

//////
// Factory.
//////

// New creates a new Ollama provider.
func New(
	options ...provider.ClientFunc,
) (*Ollama, error) {
	// Enforces IProvider interface implementation.
	var _ provider.IProvider = (*Ollama)(nil)

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

	provider := &Ollama{
		Provider: p,

		Endpoint: p.Endpoint,

		client: client,
	}

	if err := validation.Validate(provider); err != nil {
		return nil, err
	}

	singleton = provider

	return provider, nil
}

// NewDefault creates a new Ollama provider with default values.
func NewDefault(options ...provider.ClientFunc) (*Ollama, error) {
	opts := []provider.ClientFunc{
		provider.WithEndpoint(config.Get().OllamaEndpoint),
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
