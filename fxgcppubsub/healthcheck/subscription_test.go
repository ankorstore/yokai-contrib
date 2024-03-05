package healthcheck_test

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	fxgcppubsubhealthcheck "github.com/ankorstore/yokai-contrib/fxgcppubsub/healthcheck"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestWithExistingSubscriptions(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	var client *pubsub.Client
	var conf *config.Config

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxgcppubsub.FxGcpPubSubModule,
		fxconfig.FxConfigModule,
		fx.Populate(&client, &conf),
	).RequireStart()

	topic, err := client.CreateTopic(context.Background(), "topic1")
	assert.NoError(t, err)

	for _, subscription := range conf.GetStringSlice("modules.gcppubsub.healthcheck.subscriptions") {
		_, err = client.CreateSubscription(context.Background(), subscription, pubsub.SubscriptionConfig{
			Topic: topic,
		})
		assert.NoError(t, err)
	}

	p := fxgcppubsubhealthcheck.NewGcpPubSubSubscriptionsProbe(conf, client)
	assert.Equal(t, fxgcppubsubhealthcheck.SubscriptionsProbeName, p.Name())
	checkResult := p.Check(context.Background())

	app.RequireStop()

	assert.True(t, checkResult.Success)
	assert.Equal(
		t,
		"subscription subscription1 exists, subscription subscription2 exists, subscription subscription3 exists",
		checkResult.Message,
	)
}

func TestWithMissingSubscriptions(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	var client *pubsub.Client
	var conf *config.Config

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxgcppubsub.FxGcpPubSubModule,
		fxconfig.FxConfigModule,
		fx.Populate(&client, &conf),
	).RequireStart()

	topic, err := client.CreateTopic(context.Background(), "topic1")
	assert.NoError(t, err)

	for _, subscription := range conf.GetStringSlice("modules.gcppubsub.healthcheck.subscriptions") {
		if subscription != "subscription2" {
			_, err = client.CreateSubscription(context.Background(), subscription, pubsub.SubscriptionConfig{
				Topic: topic,
			})
			assert.NoError(t, err)
		}
	}

	p := fxgcppubsubhealthcheck.NewGcpPubSubSubscriptionsProbe(conf, client)
	assert.Equal(t, fxgcppubsubhealthcheck.SubscriptionsProbeName, p.Name())
	checkResult := p.Check(context.Background())

	app.RequireStop()

	assert.False(t, checkResult.Success)
	assert.Equal(
		t,
		"subscription subscription1 exists, subscription subscription2 does not exist, subscription subscription3 exists",
		checkResult.Message,
	)
}

func TestWithEmptySubscriptions(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	// empty the subscriptions list
	t.Setenv("MODULES_GCPPUBSUB_HEALTHCHECK_SUBSCRIPTIONS", " ")

	var client *pubsub.Client
	var conf *config.Config

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxgcppubsub.FxGcpPubSubModule,
		fxconfig.FxConfigModule,
		fx.Populate(&client, &conf),
	).RequireStart()

	p := fxgcppubsubhealthcheck.NewGcpPubSubSubscriptionsProbe(conf, client)
	assert.Equal(t, fxgcppubsubhealthcheck.SubscriptionsProbeName, p.Name())
	checkResult := p.Check(context.Background())

	app.RequireStop()

	assert.True(t, checkResult.Success)
}

func TestWithFailingSubscriptions(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	var client *pubsub.Client
	var conf *config.Config

	fxtest.New(
		t,
		fx.NopLogger,
		fxgcppubsub.FxGcpPubSubModule,
		fxconfig.FxConfigModule,
		fx.Populate(&client, &conf),
	).RequireStart().RequireStop()

	p := fxgcppubsubhealthcheck.NewGcpPubSubSubscriptionsProbe(conf, client)
	assert.Equal(t, fxgcppubsubhealthcheck.SubscriptionsProbeName, p.Name())
	checkResult := p.Check(context.Background())

	assert.False(t, checkResult.Success)
	assert.Contains(t, checkResult.Message, "subscription subscription1 error: rpc error")
	assert.Contains(t, checkResult.Message, "subscription subscription2 error: rpc error")
	assert.Contains(t, checkResult.Message, "subscription subscription3 error: rpc error")
}
