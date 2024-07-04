package subscription_test

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/subscription"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestDefaultSubscriptionRegistry(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	var registry subscription.SubscriptionRegistry
	var client *pubsub.Client

	ctx := context.Background()

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxgcppubsub.FxGcpPubSubModule,
		fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
		fxgcppubsub.PrepareTopicAndSubscription(fxgcppubsub.PrepareTopicAndSubscriptionParams{
			TopicID:        "test-topic",
			SubscriptionID: "test-subscription",
		}),
		fx.Populate(&registry, &client),
	).RequireStart().RequireStop()

	t.Run("registry empty by default", func(t *testing.T) {
		assert.Len(t, registry.All(), 0)

		assert.False(t, registry.Has("test-subscription"))

		_, err := registry.Get("test-subscription")
		assert.Error(t, err)
		assert.Equal(t, "cannot find subscription test-subscription", err.Error())
	})

	t.Run("registry lifecycle", func(t *testing.T) {
		baseSub := client.Subscription("test-subscription")

		sub := subscription.NewSubscription(
			codec.NewDefaultCodec(pubsub.SchemaTypeUnspecified, pubsub.EncodingUnspecified, ""),
			baseSub,
		)

		registry.Add(sub)

		assert.Len(t, registry.All(), 1)

		assert.True(t, registry.Has("test-subscription"))

		registeredSub, err := registry.Get("test-subscription")
		assert.NoError(t, err)
		assert.Equal(t, sub, registeredSub)
	})
}
