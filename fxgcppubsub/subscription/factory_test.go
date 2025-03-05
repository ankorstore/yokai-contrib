package subscription_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/subscription"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestDefaultSubscriptionFactory(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	var factory subscription.SubscriptionFactory

	ctx := context.Background()

	t.Run("subscription creation", func(t *testing.T) {
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
			fx.Populate(&factory),
		).RequireStart().RequireStop()

		sub, err := factory.Create(ctx, "test-subscription")
		assert.NoError(t, err)

		assert.Equal(t, "test-subscription", sub.BaseSubscription().ID())
		assert.Equal(t, "projects/test-project/subscriptions/test-subscription", sub.BaseSubscription().String())
	})

	t.Run("subscription creation with schema", func(t *testing.T) {
		fxtest.New(
			t,
			fx.NopLogger,
			fxconfig.FxConfigModule,
			fxlog.FxLogModule,
			fxgcppubsub.FxGcpPubSubModule,
			fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
			fxgcppubsub.PrepareTopicAndSubscriptionWithSchema(fxgcppubsub.PrepareTopicAndSubscriptionWithSchemaParams{
				TopicID:        "test-topic",
				SubscriptionID: "test-subscription",
				SchemaID:       "test-schema",
			}),
			fx.Populate(&factory),
		).RequireStart().RequireStop()

		sub, err := factory.Create(ctx, "test-subscription")
		assert.NoError(t, err)

		assert.Equal(t, "test-subscription", sub.BaseSubscription().ID())
		assert.Equal(t, "projects/test-project/subscriptions/test-subscription", sub.BaseSubscription().String())
	})

	t.Run("subscription creation error", func(t *testing.T) {
		fxtest.New(
			t,
			fx.NopLogger,
			fxconfig.FxConfigModule,
			fxlog.FxLogModule,
			fxgcppubsub.FxGcpPubSubModule,
			fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
			fx.Populate(&factory),
		).RequireStart().RequireStop()

		sub, err := factory.Create(ctx, "test-subscription")
		assert.Nil(t, sub)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot get subscription test-subscription configuration")
	})
}
