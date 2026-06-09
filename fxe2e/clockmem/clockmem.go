// Package clockmem pins the application clock to a case's fixed time for
// end-to-end tests.
//
// When a case declares a fixed wall-clock time via the request's "now" field
// (RFC 3339), wiring clockmem.FxOption into the Boot replaces the app's
// clockwork.Clock with a fake clock anchored at that instant, so time-dependent
// logic (timestamps, expiries, "is it stale yet") runs deterministically. When
// the case declares no "now", it is a no-op and the real clock is used.
//
// Wire it into a Boot alongside the other helpers — no per-case code:
//
//	internal.RunE2ETest(tb,
//	    clockmem.FxOption(fix),
//	    // ... other options
//	)
//
// It assumes the app resolves time through yokai's fxclock (clockwork.Clock),
// which is the Ankorstore convention.
package clockmem

import (
	"github.com/ankorstore/yokai-contrib/fxe2e/e2e"
	"github.com/jonboulle/clockwork"
	"go.uber.org/fx"
)

// FxOption returns an fx.Option that pins clockwork.Clock to the fixture's fixed
// time, or a no-op option when the case declares no "now".
func FxOption(fix *e2e.Fixture) fx.Option {
	now, ok := fix.FixedTime()
	if !ok {
		return fx.Options()
	}

	return fx.Decorate(func(clockwork.Clock) clockwork.Clock {
		return clockwork.NewFakeClockAt(now)
	})
}
