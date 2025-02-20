package fxjsonapi

import (
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxhttpserver"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "jsonapi"

// FxJSONAPIModule is the [Fx] JSON API module.
//
// [Fx]: https://github.com/uber-go/fx
var FxJSONAPIModule = fx.Module(
	ModuleName,
	fxhttpserver.AsErrorHandler(NewErrorHandler),
	fx.Provide(fx.Annotate(ProvideProcessor, fx.As(new(Processor)))),
)

// ProvideProcessorParam allows injection of the required dependencies in ProvideProcessor.
type ProvideProcessorParam struct {
	fx.In
	Config *config.Config
}

// ProvideProcessor provides a new DefaultProcessor instance.
func ProvideProcessor(p ProvideProcessorParam) *DefaultProcessor {
	return NewDefaultProcessor(p.Config)
}
