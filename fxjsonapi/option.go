package fxjsonapi

import "github.com/ankorstore/yokai/config"

type Options struct {
	Metadata map[string]any
	Included bool
	Log      bool
	Trace    bool
}

func DefaultProcessorOptions(config *config.Config) Options {
	return Options{
		Metadata: make(map[string]any),
		Included: true,
		Log:      config.GetBool("modules.jsonapi.log.enabled"),
		Trace:    config.GetBool("modules.jsonapi.trace.enabled"),
	}
}

type ProcessorOption func(o *Options)

func WithMetadata(m map[string]any) ProcessorOption {
	return func(o *Options) {
		o.Metadata = m
	}
}

func WithoutIncluded() ProcessorOption {
	return func(o *Options) {
		o.Included = false
	}
}

func WithLog() ProcessorOption {
	return func(o *Options) {
		o.Log = true
	}
}

func WithTrace() ProcessorOption {
	return func(o *Options) {
		o.Trace = true
	}
}
