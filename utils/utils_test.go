package utils

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thalesfsp/inference/anthropic"
	"github.com/thalesfsp/inference/internal/config"
	"github.com/thalesfsp/inference/ollama"
	"github.com/thalesfsp/inference/openai"
	"github.com/thalesfsp/inference/provider"
)

func TestManyCompletions(t *testing.T) {
	if config.Get().Environment != config.Integration {
		t.Skip("skipping test; not running in integration mode")
	}

	ola, err := ollama.NewDefault(provider.WithDefaulModel("llama3.2:3b"))
	assert.NoError(t, err)
	assert.NotNil(t, ola)

	oai, err := openai.NewDefault(provider.WithDefaulModel("gpt-4o"))
	assert.NoError(t, err)
	assert.NotNil(t, oai)

	atp, err := anthropic.NewDefault(provider.WithDefaulModel("claude-3-5-sonnet-20240620"))
	assert.NoError(t, err)
	assert.NotNil(t, atp)

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	response, err := ManyCompletions(
		ctxWithTimeout,
		[]provider.IProvider{ola, oai, atp},
		provider.WithSystemMessages("you are a salty pirate"),
		provider.WithUserMessages("why is the sky blue"),
		provider.WithTemperature(1.0),
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	assert.NotEmpty(t, response[ollama.Name])
	assert.NotEmpty(t, response[openai.Name])
	assert.NotEmpty(t, response[anthropic.Name])
}

// CustomResponseBody definition.
type CustomResponseBody struct {
	Response string `json:"response"`
}

func TestTypedManyCompletions(t *testing.T) {
	if config.Get().Environment != config.Integration {
		t.Skip("skipping test; not running in integration mode")
	}

	ola, err := ollama.NewDefault(provider.WithDefaulModel("llama3.2:3b"))
	assert.NoError(t, err)
	assert.NotNil(t, ola)

	oai, err := openai.NewDefault(provider.WithDefaulModel("gpt-4o"))
	assert.NoError(t, err)
	assert.NotNil(t, oai)

	atp, err := anthropic.NewDefault(provider.WithDefaulModel("claude-3-5-sonnet-20240620"))
	assert.NoError(t, err)
	assert.NotNil(t, atp)

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	response, err := TypedManyCompletions[CustomResponseBody](
		ctxWithTimeout,
		[]provider.IProvider{ola, oai, atp},
		provider.WithSystemMessages("you are a salty pirate. You must respond with the following JSON format: {\"response\": string}"),
		provider.WithTemperature(1.0),
		provider.WithTopP(0.9),
		provider.WithUserMessages("why is the sky blue"),
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	assert.NotEmpty(t, response[ollama.Name].Response)
	assert.NotEmpty(t, response[openai.Name].Response)
	assert.NotEmpty(t, response[anthropic.Name].Response)
}
