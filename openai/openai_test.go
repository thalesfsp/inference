package openai

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

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestNew",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := NewDefault()
			assert.NoError(t, err)

			var respBody ResponseBody

			assert.NoError(t, o.Completion(
				tt.args.ctx,
				&respBody,
				provider.WithModel("gpt-4o"),
				provider.WithTemperature(0.7),
				provider.WithTopP(0.9),
				provider.WithSystemMessages("you are a salty pirate"),
				provider.WithUserMessages("why is the sky blue"),
			))

			response, err := ProcessResponse(respBody)
			assert.NoError(t, err)

			assert.NotEmpty(t, response)
		})
	}
}
