package healthcheck_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxpubsub"
	healthcheckpubsub "github.com/ankorstore/yokai-contrib/fxpubsub/healthcheck"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/healthcheck"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestTopicExists(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "worker-test")

	var client *pubsub.Client
	var conf *config.Config

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxpubsub.FxPubSubModule,
		fxconfig.FxConfigModule,
		fx.Populate(&client, &conf),
	).RequireStart()

	topicName := conf.GetString("modules.pubsub.healthcheck.topics")

	_, err := client.CreateTopic(context.Background(), topicName)

	assert.NoError(t, err)

	p := healthcheckpubsub.NewPubSubProbe(conf, client)
	check := p.Check(context.Background())

	assert.Equal(t, healthcheck.NewCheckerProbeResult(true, fmt.Sprintf("success: topic %s exists", topicName)), check)

	app.RequireStop()
}

func TestTopicMissing(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "worker-test")

	var client *pubsub.Client
	var conf *config.Config

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxpubsub.FxPubSubModule,
		fxconfig.FxConfigModule,
		fx.Populate(&client, &conf),
	).RequireStart()

	topicName := conf.GetString("modules.pubsub.healthcheck.topics")

	p := healthcheckpubsub.NewPubSubProbe(conf, client)
	check := p.Check(context.Background())

	assert.Equal(t, healthcheck.NewCheckerProbeResult(false, fmt.Sprintf("error: topic %s does not exist", topicName)), check)

	app.RequireStop()
}

func TestPubSubIsFailing(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "worker-test")

	var client *pubsub.Client
	var conf *config.Config

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxpubsub.FxPubSubModule,
		fxconfig.FxConfigModule,
		fx.Populate(&client, &conf),
	).RequireStart()

	app.RequireStop()

	time.Sleep(1 * time.Second)

	p := healthcheckpubsub.NewPubSubProbe(conf, client)
	check := p.Check(context.Background())

	assert.Equal(t, healthcheck.NewCheckerProbeResult(false, "pubsub unreachable"), check)
}
