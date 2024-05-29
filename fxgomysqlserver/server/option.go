package server

import (
	"io"

	"github.com/ankorstore/yokai-contrib/fxgomysqlserver/config"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// ServerOptions are options for the [GoMySQLServerFactory] implementations.
type ServerOptions struct {
	Config    *config.GoMySQLServerConfig
	LogOutput io.Writer
	Tracer    trace.Tracer
}

// DefaultGoMySQLServerOptions are the default options used in the [DefaultGoMySQLServerFactory].
func DefaultGoMySQLServerOptions() ServerOptions {
	return ServerOptions{
		Config:    config.NewGoMySQLServerConfig(),
		LogOutput: io.Discard,
		Tracer:    noop.NewTracerProvider().Tracer("noop"),
	}
}

// GoMySQLServerOption are functional options for the [GoMySQLServerFactory] implementations.
type GoMySQLServerOption func(o *ServerOptions)

// WithConfig is used to specify the [config.GoMySQLServerConfig] to use as server config.
func WithConfig(config *config.GoMySQLServerConfig) GoMySQLServerOption {
	return func(o *ServerOptions) {
		o.Config = config
	}
}

// WithLogOutput is used to specify the [io.Writer] to use as server logger output.
func WithLogOutput(output io.Writer) GoMySQLServerOption {
	return func(o *ServerOptions) {
		o.LogOutput = output
	}
}

// WithTracer is used to specify the [trace.Tracer] to use as server tracer.
func WithTracer(tracer trace.Tracer) GoMySQLServerOption {
	return func(o *ServerOptions) {
		o.Tracer = tracer
	}
}
