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

func TestGcpPubSubSubscriptionsProbe(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	var config *config.Config
	var client *pubsub.Client

	ctx := context.Background()

	t.Run("probe name", func(t *testing.T) {
		probe := &healthcheck.GcpPubSubSubscriptionsProbe{}

		assert.Equal(t, "gcppubsub-subscriptions", probe.Name())
	})

	t.Run("probe success when subscription exist", func(t *testing.T) {
		fxtest.New(
			t,
			fx.NopLogger,
			fxconfig.FxConfigModule,
			fxlog.FxLogModule,
			fxgcppubsub.FxGcpPubSubModule,
			fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
			fxgcppubsub.PrepareTopicAndSubscription(fxgcppubsub.PrepareTopicAndSubscriptionParams{
				TopicID:        "test-topic",
				SubscriptionID: "test-subscription",
			}),
			fx.Populate(&config, &client),
		).RequireStart().RequireStop()

		probe := healthcheck.NewGcpPubSubSubscriptionsProbe(config, client)

		res := probe.Check(ctx)
		assert.True(t, res.Success)
		assert.Equal(t, "subscription test-subscription exists", res.Message)
	})

	t.Run("probe failure when subscription does not exist", func(t *testing.T) {
		fxtest.New(
			t,
			fx.NopLogger,
			fxconfig.FxConfigModule,
			fxlog.FxLogModule,
			fxgcppubsub.FxGcpPubSubModule,
			fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
			fx.Populate(&config, &client),
		).RequireStart().RequireStop()

		probe := healthcheck.NewGcpPubSubSubscriptionsProbe(config, client)

		res := probe.Check(ctx)
		assert.False(t, res.Success)
		assert.Equal(t, "subscription test-subscription does not exist", res.Message)
	})
}
