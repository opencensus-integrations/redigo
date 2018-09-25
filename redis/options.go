package redis

import (
	"go.opencensus.io/trace"
)

// TraceOption allows for managing redigo configuration using functional options.
// Copy from https://github.com/opencensus-integrations/ocsql/blob/master/options.go
type TraceOption func(o *TraceOptions)

type TraceOptions struct {
	// DefaultAttributes will be set to each span as default.
	DefaultAttributes []trace.Attribute
}

// WithDefaultAttributes will be set to each span as default.
func WithDefaultAttributes(attrs ...trace.Attribute) TraceOption {
	return func(o *TraceOptions) {
		o.DefaultAttributes = attrs
	}
}
