package ollama

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thalesfsp/inference/internal/config"
	"github.com/thalesfsp/inference/provider"
)

func TestNew(t *testing.T) {
	if config.Get().Environment != config.Integration {
		t.Skip("skipping test; not running in integration mode")
	}

	// Model to be used in the test.
	model := "llama3.2:3b"

	// CustomResponseBody definition.
	type CustomResponseBody struct {
		Response string `json:"response"`
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name             string
		args             args
		model            string
		systemMessage    string
		userMessage      string
		withResponseBody bool
	}{
		{
			name: "TestNew",
			args: args{
				ctx: context.Background(),
			},
			model:         model,
			systemMessage: "you are a salty pirate",
			userMessage:   "why is the sky blue",
		},
		{
			name: "TestNew",
			args: args{
				ctx: context.Background(),
			},
			model:            model,
			systemMessage:    "you are a salty pirate. You must respond with the following JSON format: {\"response\": string}",
			userMessage:      "why is the sky blue",
			withResponseBody: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := NewDefault()
			assert.NoError(t, err)
			assert.NotNil(t, o)

			options := []provider.Func{
				provider.WithModel(tt.model),
				provider.WithTemperature(1.0),
				provider.WithTopK(40),
				provider.WithTopP(0.9),
				provider.WithSystemMessages(tt.systemMessage),
				provider.WithUserMessages(tt.userMessage),
			}

			var customResponseBody CustomResponseBody

			if tt.withResponseBody {
				options = append(
					options,
					provider.WithResponseBody(&customResponseBody),
				)
			}

			response, err1 := o.Completion(
				tt.args.ctx,
				options...,
			)

			assert.NoError(t, err1)
			assert.NotEmpty(t, response)

			if tt.withResponseBody {
				assert.NotEmpty(t, customResponseBody.Response)
			}
		})
	}
}
