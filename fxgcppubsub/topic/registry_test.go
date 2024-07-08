package topic_test

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/topic"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestDefaultTopicRegistry(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	var registry topic.TopicRegistry
	var client *pubsub.Client

	ctx := context.Background()

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxgcppubsub.FxGcpPubSubModule,
		fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
		fxgcppubsub.PrepareTopic(fxgcppubsub.PrepareTopicParams{
			TopicID: "test-topic",
		}),
		fx.Populate(&registry, &client),
	).RequireStart().RequireStop()

	t.Run("registry empty by default", func(t *testing.T) {
		assert.Len(t, registry.All(), 0)

		assert.False(t, registry.Has("test-topic"))

		_, err := registry.Get("test-topic")
		assert.Error(t, err)
		assert.Equal(t, "cannot find topic test-topic", err.Error())
	})

	t.Run("registry lifecycle", func(t *testing.T) {
		baseSub := client.Topic("test-topic")

		top := topic.NewTopic(
			codec.NewRawCodec(),
			baseSub,
		)

		registry.Add(top)

		assert.Len(t, registry.All(), 1)

		assert.True(t, registry.Has("test-topic"))

		registeredTop, err := registry.Get("test-topic")
		assert.NoError(t, err)
		assert.Equal(t, top, registeredTop)
	})
}
