package fxgcppubsub

import (
	"cloud.google.com/go/pubsub/pstest"
	"go.uber.org/fx"
)

// Reactor is the interface for pub/sub test server reactors.
type Reactor interface {
	FuncNames() []string
	pstest.Reactor
}

// AsPubSubTestServerReactor registers a [Reactor] into Fx.
func AsPubSubTestServerReactor(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(Reactor)),
			fx.ResultTags(`group:"gcppubsub-reactors"`),
		),
	)
}

// AsPubSubTestServerReactors registers a list of [Reactor] into Fx.
func AsPubSubTestServerReactors(constructors ...any) fx.Option {
	options := []fx.Option{}

	for _, constructor := range constructors {
		options = append(options, AsPubSubTestServerReactor(constructor))
	}

	return fx.Options(options...)
}
