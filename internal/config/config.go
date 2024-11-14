package config

import (
	"log"
	"sync"
	"time"

	"github.com/thalesfsp/configurer/cmd"
	"github.com/thalesfsp/configurer/util"
)

//////
// Consts, vars, and types.
//////

var config string

// Singleton.
var (
	once            sync.Once
	singletonConfig *Config
)

// Config is the application configuration.
//
// NOTE: Mark sensitive fields with `json:"-"` to avoid them being dumped.
//
//nolint:lll
type Config struct {
	// NOTE: For integration testing, set `environment` to `integration`.
	Environment string `default:"testing" env:"ENVIRONMENT" json:"environment" validate:"required,oneof=testing integration"`

	// Anthropic.
	AnthropicEndpoint string `default:"https://api.anthropic.com/v1/messages" env:"ANTHROPIC_ENDPOINT" json:"anthropicBaseURL"   validate:"omitempty,gt=0"`
	AnthropicToken    string `env:"ANTHROPIC_API_KEY"                         json:"-"                 validate:"omitempty,gt=0"`

	// HuggingFace.
	HuggingFaceEndpoint string `default:"https://api-inference.huggingface.co/v1/chat/completions" env:"HUGGINGFACE_ENDPOINT" json:"huggingFaceBaseURL" validate:"omitempty,gt=0"`
	HuggingFaceToken    string `env:"HUGGINGFACE_API_KEY"                                          json:"-"                   validate:"omitempty,gt=0"`

	// OpenAI.
	OpenAIEndpoint string `default:"https://api.openai.com/v1/chat/completions" env:"OPEN_AI_ENDPOINT" json:"openAIBaseURL"      validate:"omitempty,gt=0"`
	OpenAIToken    string `env:"OPENAI_API_KEY"                                 json:"-"               validate:"omitempty,gt=0"`

	// Ollama.
	OllamaEndpoint string `default:"http://localhost:11434/api/chat" env:"OLLAMA_ENDPOINT" json:"ollamaBaseURL" validate:"omitempty,gt=0"`

	//////
	// Common timeouts.
	//////

	TimeoutLong   time.Duration `default:"30s" env:"TIMEOUT_LONG"  json:"timeoutLong"  validate:"omitempty,gt=0"`
	TimeoutMedium time.Duration `default:"10s" env:"TIMEOUT_SHORT" json:"timeoutShort" validate:"omitempty,gt=0"`
	TimeoutShort  time.Duration `default:"3s"  env:"TIMEOUT_SHORT" json:"timeoutShort" validate:"omitempty,gt=0"`
}

//////
// Exported functionalities.
//////

// Get returns a setup config.
func Get() *Config {
	once.Do(func() {
		singletonConfig = Reload()
	})

	return singletonConfig
}

// Reload reloads the configuration.
func Reload() *Config {
	if _, err := cmd.LoadFromText(false, false, "env", config); err != nil {
		log.Fatalln("Failed to load env var text", err)
	}

	// Create a virtual file to be passed to dotenv.New(file string).

	conf := &Config{}

	// Configuration was already exported to env var. API binary is wrapped
	// by the configuration cli.
	if err := util.Process(conf); err != nil {
		log.Fatalln("Failed to process configuration", err)
	}

	// Update the global configuration.
	singletonConfig = conf

	return conf
}
