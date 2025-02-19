package fxjsonapi

import (
	"github.com/ankorstore/yokai/config"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "jsonapi"

var FxJSONAPIModule = fx.Module(
	ModuleName,
	fx.Provide(
		fx.Annotate(ProvideProcessor, fx.As(new(Processor))),
	),
)

type ProvideProcessorParam struct {
	fx.In
	Config *config.Config
}

func ProvideProcessor(p ProvideProcessorParam) *DefaultProcessor {
	return NewDefaultProcessor(p.Config)
}
