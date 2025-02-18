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
		ProvideProcessor,
	),
)

type ProvideProcessorParam struct {
	fx.In
	Config *config.Config
}

func ProvideProcessor(p ProvideProcessorParam) Processor {
	return NewDefaultProcessor(p.Config)
}
