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

func TestWithExistingTopics(t *testing.T) {
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

	for _, topic := range conf.GetStringSlice("modules.gcppubsub.healthcheck.topics") {
		_, err := client.CreateTopic(context.Background(), topic)
		assert.NoError(t, err)
	}

	p := fxgcppubsubhealthcheck.NewGcpPubSubTopicsProbe(conf, client)
	assert.Equal(t, fxgcppubsubhealthcheck.TopicsProbeName, p.Name())
	checkResult := p.Check(context.Background())

	app.RequireStop()

	assert.True(t, checkResult.Success)
	assert.Equal(t, "topic topic1 exists, topic topic2 exists, topic topic3 exists", checkResult.Message)
}

func TestWithMissingTopics(t *testing.T) {
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

	for _, topic := range conf.GetStringSlice("modules.gcppubsub.healthcheck.topics") {
		if topic != "topic2" {
			_, err := client.CreateTopic(context.Background(), topic)
			assert.NoError(t, err)
		}
	}

	p := fxgcppubsubhealthcheck.NewGcpPubSubTopicsProbe(conf, client)
	assert.Equal(t, fxgcppubsubhealthcheck.TopicsProbeName, p.Name())
	checkResult := p.Check(context.Background())

	app.RequireStop()

	assert.False(t, checkResult.Success)
	assert.Equal(t, "topic topic1 exists, topic topic2 does not exist, topic topic3 exists", checkResult.Message)
}

func TestWithEmptyTopics(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	// empty the topics list
	t.Setenv("MODULES_GCPPUBSUB_HEALTHCHECK_TOPICS", " ")

	var client *pubsub.Client
	var conf *config.Config

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxgcppubsub.FxGcpPubSubModule,
		fxconfig.FxConfigModule,
		fx.Populate(&client, &conf),
	).RequireStart()

	p := fxgcppubsubhealthcheck.NewGcpPubSubTopicsProbe(conf, client)
	assert.Equal(t, fxgcppubsubhealthcheck.TopicsProbeName, p.Name())
	checkResult := p.Check(context.Background())

	app.RequireStop()

	assert.True(t, checkResult.Success)
}

func TestWithFailingTopics(t *testing.T) {
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

	p := fxgcppubsubhealthcheck.NewGcpPubSubTopicsProbe(conf, client)
	assert.Equal(t, fxgcppubsubhealthcheck.TopicsProbeName, p.Name())
	checkResult := p.Check(context.Background())

	assert.False(t, checkResult.Success)
	assert.Contains(t, checkResult.Message, "topic topic1 error: rpc error")
	assert.Contains(t, checkResult.Message, "topic topic2 error: rpc error")
	assert.Contains(t, checkResult.Message, "topic topic3 error: rpc error")
}
