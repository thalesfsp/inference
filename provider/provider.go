package provider

import (
	"expvar"

	"github.com/thalesfsp/inference/internal/metrics"
	"github.com/thalesfsp/status"
	"github.com/thalesfsp/sypl"
	"github.com/thalesfsp/sypl/level"
	"github.com/thalesfsp/validation"
)

//////
// Vars, consts, and types.
//////

// Type is the type of the entity regarding the framework. It is used to for
// example, to identify the entity in the logs, metrics, and for tracing.
const (
	Type = "provider"
)

// Provider definition.
type Provider struct {
	// Logger.
	Logger sypl.ISypl `json:"-" validate:"required"`

	// Name of the provider type.
	Name string `json:"name" validate:"required,lowercase,gte=1"`

	// Metrics.
	counterCompletion       *expvar.Int `json:"-" validate:"required,gte=0"`
	counterCompletionFailed *expvar.Int `json:"-" validate:"required,gte=0"`

	// A provider may have the following...
	// Endpoint to reach the provider.
	Endpoint string `json:"endpoint,omitempty"`

	// DefaultModel default model to be used.
	DefaultModel string `json:"model,omitempty"`

	// Token to authenticate against the provider.
	Token string `json:"-"`
}

//////
// Implements the IMeta interface.
//////

// GetLogger returns the logger.
func (s *Provider) GetLogger() sypl.ISypl {
	return s.Logger
}

// GetName returns the provider name.
func (s *Provider) GetName() string {
	return s.Name
}

// GetType returns its type.
func (s *Provider) GetType() string {
	return Type
}

//////
// Implements the IMetrics interface.
//////

// GetCounterCounted returns the metric.
func (s *Provider) GetCounterCounted() *expvar.Int {
	return s.counterCompletion
}

// GetCounterCompletion returns the completion metric.
func (s *Provider) GetCounterCompletion() *expvar.Int {
	return s.counterCompletion
}

// GetCounterCompletionFailed returns the failed completion metric.
func (s *Provider) GetCounterCompletionFailed() *expvar.Int {
	return s.counterCompletionFailed
}

//////
// Factory.
//////

// New returns a new provider.
func New(
	name string,
	options ...ClientFunc,
) (*Provider, error) {
	//////
	// Provider options processing.
	//////

	defaultProviderOptions := ClientOptions{}

	// Apply options against the default options.
	for _, option := range options {
		if err := option(&defaultProviderOptions); err != nil {
			return nil, err
		}
	}

	// Validate the default options.
	if err := validation.Validate(&defaultProviderOptions); err != nil {
		return nil, err
	}

	//////
	// Provider setup.
	//////

	// Provider's individual logger.
	logger := sypl.NewDefault(name, level.Error).SetTags(Type, name)

	a := &Provider{
		Logger: logger,
		Name:   name,

		counterCompletion:       metrics.NewIntCounter(Type, name, "completion"),
		counterCompletionFailed: metrics.NewIntCounter(Type, name, "completion"+"."+status.Failed.String()),

		//////
		// A provider may have the following...
		//////

		Endpoint:     defaultProviderOptions.Endpoint,
		DefaultModel: defaultProviderOptions.Model,
		Token:        defaultProviderOptions.Token,
	}

	// Validate the provider.
	if err := validation.Validate(a); err != nil {
		return nil, err
	}

	// Notify the creation.
	a.GetLogger().Debugln(status.Created.String())

	return a, nil
}
