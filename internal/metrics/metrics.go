package metrics

import (
	"expvar"
	"fmt"
	"time"
)

// DefaultMetricCounterLabel is the default label for the metric counter.
const DefaultMetricCounterLabel = "counter"

// NewInt creates and initializes a new expvar.Int.
func NewInt(
	entityType,
	entityName,
	status,
	metricType string,
) *expvar.Int {
	timeInUnix := time.Now().Unix()

	counter := expvar.NewInt(fmt.Sprintf(
		"%s.%s.%d.%s.%s",
		entityType,
		entityName,
		timeInUnix,
		status,
		metricType,
	),
	)

	counter.Set(0)

	return counter
}

// NewIntCounter creates and initializes a new expvar.Int counter.
func NewIntCounter(entityType, entityName, name string) *expvar.Int {
	return NewInt(
		entityType,
		entityName,
		name,
		DefaultMetricCounterLabel,
	)
}
