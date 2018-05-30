package observability

import (
	"context"
	"time"

	"go.opencensus.io/tag"
)

func TagKeyValuesIntoContext(ctx context.Context, key tag.Key, values ...string) (context.Context, error) {
	insertions := make([]tag.Mutator, len(values))
	for i, value := range values {
		insertions[i] = tag.Insert(key, value)
	}
	return tag.New(ctx, insertions...)
}

func SinceInMilliseconds(startTime time.Time) float64 {
	return float64(time.Since(startTime).Nanoseconds()) * 1e6
}
