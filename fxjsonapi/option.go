package fxjsonapi

import "github.com/ankorstore/yokai/config"

// Options are options for the [Processor].
type Options struct {
	Metadata map[string]any
	Included bool
	Log      bool
	Trace    bool
}

// DefaultProcessorOptions are the default [Processor] options.
func DefaultProcessorOptions(config *config.Config) Options {
	return Options{
		Metadata: make(map[string]any),
		Included: true,
		Log:      config.GetBool("modules.jsonapi.log.enabled"),
		Trace:    config.GetBool("modules.jsonapi.trace.enabled"),
	}
}

// ProcessorOption are functional options for the [Processor].
type ProcessorOption func(o *Options)

// WithMetadata is used to add metadata to the json api representation.
func WithMetadata(m map[string]any) ProcessorOption {
	return func(o *Options) {
		o.Metadata = m
	}
}

// WithIncluded is used to add included to the json api representation.
func WithIncluded(i bool) ProcessorOption {
	return func(o *Options) {
		o.Included = i
	}
}

// WithLog is used to add logging.
func WithLog(l bool) ProcessorOption {
	return func(o *Options) {
		o.Log = l
	}
}

// WithTrace is used to add tracing.
func WithTrace(t bool) ProcessorOption {
	return func(o *Options) {
		o.Trace = t
	}
}
