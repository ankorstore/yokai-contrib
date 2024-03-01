package fxpubsub_test

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
	fxpubsub "github.com/ankorstore/yokai-contrib/fxpubsub"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestFxPubSubModule(t *testing.T) {
	app := fxtest.New(
		t,
		fx.NopLogger,
		fxpubsub.FxPubSubModule,
		fxconfig.FxConfigModule,
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to start the Fx application")

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to stop the Fx application")
}

func TestNewFxPubSubForTestClientWithoutProjectID(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config/")

	var conf *config.Config
	var client *pubsub.Client

	testApp := fxtest.New(
		t,
		fx.NopLogger,
		fxpubsub.FxPubSubModule,
		fxconfig.FxConfigModule,
		fx.Populate(&conf, &client),
	)

	err := testApp.Start(context.Background())
	assert.Error(t, err, "failed to create test pubsub client: pubsub: projectID string is empty")
	assert.Nil(t, client)
}

func TestNewFxPubSubForTestClient(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("GCP_PROJECT_ID", "pubsub-test")

	var conf *config.Config
	var client *pubsub.Client

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxpubsub.FxPubSubModule,
		fxconfig.FxConfigModule,
		fx.Populate(&conf, &client),
	)

	app.RequireStart()
	assert.NoError(t, app.Err(), "failed to create pubsub.Client")
	assert.NotNil(t, client)
	assert.Equal(t, "pubsub-test", client.Project())

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to close pubsub.Client")
}

func TestNewFxPubSubWithoutProjectID(t *testing.T) {
	t.Setenv("APP_ENV", "dev")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var conf *config.Config
	var client *pubsub.Client

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxpubsub.FxPubSubModule,
		fxconfig.FxConfigModule,
		fx.Populate(&conf, &client),
	)

	app.RequireStart()
	assert.Error(t, app.Err(), "failed to create pubsub client: pubsub: projectID string is empty")
}

func TestNewFxPubSubClient(t *testing.T) {
	t.Setenv("APP_ENV", "dev")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("GCP_PROJECT_ID", "pubsub")

	var conf *config.Config
	var client *pubsub.Client

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxpubsub.FxPubSubModule,
		fxconfig.FxConfigModule,
		fx.Populate(&conf, &client),
	).RequireStart()
	assert.NoError(t, app.Err(), "failed to create pubsub.Client")
	assert.NotNil(t, client)
	assert.Equal(t, "pubsub", client.Project())

	app.RequireStop()
	assert.NoError(t, app.Err(), "failed to close pubsub.Client")
}
