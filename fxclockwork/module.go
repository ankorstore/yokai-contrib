package fxclockwork

import (
	"github.com/ankorstore/yokai/config"
	"github.com/jonboulle/clockwork"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "clockwork"

// FxClockworkModule is the [Fx] clockwork module.
//
// [Fx]: https://github.com/uber-go/fx
var FxClockworkModule = fx.Module(
	ModuleName,
	fx.Provide(
		NewFxClockworkClock,
	),
)

// FxClockworkClockParam allows injection of the required dependencies in [NewFxClockwork].
type FxClockworkClockParam struct {
	fx.In
	Config *config.Config
}

// NewFxClockworkClock returns a new [clockwork.Clock] instance.
func NewFxClockworkClock(p FxClockworkClockParam) clockwork.Clock {
	if p.Config.IsTestEnv() {
		return clockwork.NewFakeClock()
	} else {
		return clockwork.NewRealClock()
	}
}
