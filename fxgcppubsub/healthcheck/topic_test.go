package healthcheck_test

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/healthcheck"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestGcpPubSubTopicsProbe(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	var config *config.Config
	var client *pubsub.Client

	ctx := context.Background()

	t.Run("probe name", func(t *testing.T) {
		probe := &healthcheck.GcpPubSubTopicsProbe{}

		assert.Equal(t, "gcppubsub-topics", probe.Name())
	})

	t.Run("probe success when topic exist", func(t *testing.T) {
		fxtest.New(
			t,
			fx.NopLogger,
			fxconfig.FxConfigModule,
			fxlog.FxLogModule,
			fxgcppubsub.FxGcpPubSubModule,
			fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
			fxgcppubsub.PrepareTopic(fxgcppubsub.PrepareTopicParams{
				TopicID: "test-topic",
			}),
			fx.Populate(&config, &client),
		).RequireStart().RequireStop()

		probe := healthcheck.NewGcpPubSubTopicsProbe(config, client)

		res := probe.Check(ctx)
		assert.True(t, res.Success)
		assert.Equal(t, "topic test-topic exists", res.Message)
	})

	t.Run("probe failure when topic does not exist", func(t *testing.T) {
		fxtest.New(
			t,
			fx.NopLogger,
			fxconfig.FxConfigModule,
			fxlog.FxLogModule,
			fxgcppubsub.FxGcpPubSubModule,
			fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
			fx.Populate(&config, &client),
		).RequireStart().RequireStop()

		probe := healthcheck.NewGcpPubSubTopicsProbe(config, client)

		res := probe.Check(ctx)
		assert.False(t, res.Success)
		assert.Equal(t, "topic test-topic does not exist", res.Message)
	})
}
