package utils

import (
	"context"

	"github.com/thalesfsp/concurrentloop"
	"github.com/thalesfsp/inference/provider"
)

// ManyCompletions generates one or more completions. It returns a map of the
// provider and its respective response. The response is a string. Use this
// method for simple, and straightforward completions.
//
// NOTE: Not all options are available for all providers.
//
// NOTE: Do not pass WithResponseBody, it will not be used, and the content will
// be wrong.
//
// SEE: utils_test.go for example.
func ManyCompletions(
	ctx context.Context,
	providers []provider.IProvider,
	options ...provider.Func) (map[string]string, error,
) {
	mapOfProviderToCompletion := make(map[string]string)

	if _, errs := concurrentloop.Map(
		ctx,
		providers,
		func(ctx context.Context, provider provider.IProvider) (string, error) {
			response, err := provider.Completion(ctx, options...)
			if err != nil {
				return "", err
			}

			// Add to the map.
			mapOfProviderToCompletion[provider.GetName()] = response

			return response, nil
		},
	); len(errs) > 0 {
		return nil, errs
	}

	return mapOfProviderToCompletion, nil
}

// TypedManyCompletions generates one or more completions. It returns a map
// correlating the provider and its respective response. The response `T`, is
// whatever the developer specified. This method is a more powerful, and
// flexible version of ManyCompletions. Use this in cases where the response is
// not a string, but needs to be unmarshalled (processed).
//
// NOTE: Not all options are available for all providers.
//
// NOTE: Do not pass WithResponseBody, it will not be used, and the content will
// be wrong.
//
// SEE: utils_test.go for example.
func TypedManyCompletions[T any](
	ctx context.Context,
	providers []provider.IProvider,
	options ...provider.Func) (map[string]T, error,
) {
	mapOfProviderToCompletion := make(map[string]T)

	if _, errs := concurrentloop.Map(
		ctx,
		providers,
		func(ctx context.Context, p provider.IProvider) (bool, error) {
			// Create a new instance of T.
			t := *new(T)

			// Append WithResponseBody to options.
			options = append(options, provider.WithResponseBody(&t))

			if _, err := p.Completion(ctx, options...); err != nil {
				return false, err
			}

			// Add to the map.
			mapOfProviderToCompletion[p.GetName()] = t

			return true, nil
		},
	); len(errs) > 0 {
		return nil, errs
	}

	return mapOfProviderToCompletion, nil
}
