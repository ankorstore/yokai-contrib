package fxslack_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxslack"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestFxSlackModule(t *testing.T) {
	app := fxtest.New(
		t,
		fx.NopLogger,
		fxslack.FxSlackModule,
		fxconfig.FxConfigModule,
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to start the Fx application")

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to stop the Fx application")
}

func TestFxSlackClient(t *testing.T) {
	t.Setenv("APP_ENV", "dev")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("TOKEN", "my-token")

	var conf *config.Config
	var client *slack.Client

	var roundTripperProvide = func() http.RoundTripper {
		return http.DefaultTransport
	}

	app := fx.New(
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxslack.FxSlackModule,
		fx.Populate(&conf, &client),
		fx.Provide(roundTripperProvide),
	)

	err := app.Start(context.Background())
	assert.NoError(t, err, "failed to create slack.Client")
	assert.NotNil(t, client)

	err = app.Stop(context.Background())
	assert.NoError(t, err, "failed to close slack.Client")
}
