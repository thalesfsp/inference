package provider

import (
	"context"
	"expvar"

	"github.com/thalesfsp/sypl"
)

//////
// Const, vars, and types.
//////

// IMeta definition.
type IMeta interface {
	GetClient() any

	// GetLogger returns the logger.
	GetLogger() sypl.ISypl

	// GetName returns the provider name.
	GetName() string

	// GetType returns its type.
	GetType() string
}

// IMetrics definition.
type IMetrics interface {
	// GetCounterCompletion returns the completion metric.
	GetCounterCompletion() *expvar.Int

	// GetCounterCompletionFailed returns the failed completion metric.
	GetCounterCompletionFailed() *expvar.Int
}

// IProvider defines what a provider does.
type IProvider interface {
	// Completion generates a completion using the provider API.
	Completion(ctx context.Context, v any, options ...Func) error
}
