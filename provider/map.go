package provider

import (
	"context"

	"github.com/thalesfsp/concurrentloop"
)

//////
// Vars, consts, and types.
//////

// Map is a map of strgs.
type Map map[string]IProvider

//////
// Methods.
//////

// String implements the Stringer interface.
func (m Map) String() string {
	// Iterate over the map and print the storage name.
	// Output: "1, 2, 3"
	var s string

	for k := range m {
		s += k + ", "
	}

	// Remove the last comma.
	if len(s) > 0 {
		s = s[:len(s)-2]
	}

	return s
}

// ToSlice converts Map to Slice of IProvider.
//
//nolint:prealloc
func (m Map) ToSlice() []IProvider {
	var s []IProvider

	for _, strg := range m {
		s = append(s, strg)
	}

	return s
}

//////
// 1:N Operations.
//////

// CompletionMany calls the Completion concurrently against all providers in the
// map. The response is mapped to the provider name.
func CompletionMany(
	ctx context.Context,
	m Map,
	options ...Func,
) (map[string]string, error) {
	responseMap := make(map[string]string)

	if _, errs := concurrentloop.MapM(ctx, m,
		func(ctx context.Context, providerName string, p IProvider) (string, error) {
			response, err := p.Completion(ctx, options...)
			if err != nil {
				return "", err
			}

			responseMap[providerName] = response

			return response, nil
		},
	); len(errs) > 0 {
		return nil, errs
	}

	return responseMap, nil
}
