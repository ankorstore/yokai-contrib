package fxgcppubsub_test

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestFxGcpPubSubModule(t *testing.T) {
	ctx := context.Background()

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxgcppubsub.FxGcpPubSubModule,
		fxconfig.FxConfigModule,
		fx.Supply(
			fx.Annotate(ctx, fx.As(new(context.Context))),
		),
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to start the Fx application")

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to stop the Fx application")
}

func TestFxGcpPubSubClient(t *testing.T) {
	t.Setenv("APP_ENV", "dev")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("GCP_PROJECT_ID", "project-test")
	t.Setenv("PUBSUB_EMULATOR_HOST", "localhost")

	ctx := context.Background()

	var conf *config.Config
	var client *pubsub.Client

	app := fx.New(
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxgcppubsub.FxGcpPubSubModule,
		fx.Supply(
			fx.Annotate(ctx, fx.As(new(context.Context))),
		),
		fx.Populate(&conf, &client),
	)

	err := app.Start(ctx)
	assert.NoError(t, err, "failed to create pubsub.Client")
	assert.NotNil(t, client)
	assert.Equal(t, "project-test", client.Project())

	err = app.Stop(ctx)
	assert.NoError(t, err, "failed to close pubsub.Client")
}

func TestFxGcpPubSubClientWithoutProjectId(t *testing.T) {
	t.Setenv("APP_ENV", "dev")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	ctx := context.Background()

	var conf *config.Config
	var client *pubsub.Client

	app := fx.New(
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxgcppubsub.FxGcpPubSubModule,
		fx.Supply(
			fx.Annotate(ctx, fx.As(new(context.Context))),
		),
		fx.Populate(&conf, &client),
	)

	err := app.Start(ctx)
	assert.Contains(t, err.Error(), "failed to create pubsub client: pubsub: projectID string is empty")
}
