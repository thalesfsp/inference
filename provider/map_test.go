package provider_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thalesfsp/inference/anthropic"
	"github.com/thalesfsp/inference/internal/config"
	"github.com/thalesfsp/inference/openai"
	"github.com/thalesfsp/inference/provider"
)

func TestCompletionMany(t *testing.T) {
	if config.Get().Environment != config.Integration {
		t.Skip("skipping test; not running in integration mode")
	}

	type args struct {
		ctx     context.Context
		m       provider.Map
		options []provider.Func
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := openai.NewDefault(provider.WithDefaulModel("gpt-4o"))
			assert.NoError(t, err)
			assert.NotNil(t, o)

			a, err := anthropic.NewDefault(provider.WithDefaulModel("claude-3-5-sonnet-20241022"))
			assert.NoError(t, err)
			assert.NotNil(t, a)

			options := []provider.Func{
				provider.WithTemperature(1.0),
				provider.WithSystemMessages("you are a salty pirate"),
				provider.WithUserMessages("why is the sky blue"),
			}

			response, err := provider.CompletionMany(context.Background(), map[string]provider.IProvider{
				"openai":    o,
				"anthropic": a,
			}, options...)
			assert.NoError(t, err)

			assert.NotEmpty(t, response[openai.Name])
			assert.NotEmpty(t, response[anthropic.Name])
		})
	}
}
