package fxclockwork_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxclockwork"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestFxClockworkClockModule(t *testing.T) {
	t.Setenv("APP_ENV", config.AppEnvDev)
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	runTest := func(tb testing.TB) clockwork.Clock {
		var clock clockwork.Clock

		app := fxtest.New(
			t,
			fx.NopLogger,
			fxconfig.FxConfigModule,
			fxclockwork.FxClockworkModule,
			fx.Populate(&clock),
		)

		app.RequireStart().RequireStop()
		assert.NoError(t, app.Err())

		return clock
	}

	t.Run("with real clock", func(t *testing.T) {
		clock := runTest(t)

		assert.NotNil(t, clock)
		assert.Implements(t, (*clockwork.Clock)(nil), clock)
		assert.Equal(t, "*clockwork.realClock", fmt.Sprintf("%T", clock))
	})

	t.Run("with test clock", func(t *testing.T) {
		t.Setenv("APP_ENV", config.AppEnvTest)

		clock := runTest(t)

		assert.NotNil(t, clock)
		assert.Implements(t, (*clockwork.Clock)(nil), clock)
		assert.Equal(t, "*clockwork.FakeClock", fmt.Sprintf("%T", clock))
	})

}
