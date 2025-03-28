package fxclockwork_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxclockwork"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestFxClockworkModule(t *testing.T) {
	app := fxtest.New(
		t,
		fx.NopLogger,
		fxclockwork.FxClockworkModule,
		fxconfig.FxConfigModule,
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to start the Fx application")

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to stop the Fx application")
}

func TestFxClockworkClock(t *testing.T) {
	t.Setenv("APP_ENV", config.AppEnvDev)
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var conf *config.Config
	var clock clockwork.Clock

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxclockwork.FxClockworkModule,
		fx.Populate(&conf, &clock),
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to create clockwork.Clock")
	assert.NotNil(t, clock)

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to close clockwork.Clock")
}

func TestFxClockworkTestClock(t *testing.T) {
	t.Setenv("APP_ENV", config.AppEnvTest)
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var conf *config.Config
	var clock clockwork.Clock

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxclockwork.FxClockworkModule,
		fx.Populate(&conf, &clock),
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to create test clockwork.Clock")
	assert.NotNil(t, clock)

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to close test clockwork.Clock")
}
