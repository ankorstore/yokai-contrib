package topic_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/topic"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestDefaultTopicFactory(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	var factory topic.TopicFactory

	ctx := context.Background()

	t.Run("topic creation", func(t *testing.T) {
		fxtest.New(
			t,
			fx.NopLogger,
			fxconfig.FxConfigModule,
			fxgcppubsub.FxGcpPubSubModule,
			fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
			fxgcppubsub.PrepareTopic(fxgcppubsub.PrepareTopicParams{
				TopicID: "test-topic",
			}),
			fx.Populate(&factory),
		).RequireStart().RequireStop()

		top, err := factory.Create(ctx, "test-topic")
		assert.NoError(t, err)

		assert.Equal(t, "test-topic", top.BaseTopic().ID())
	})

	t.Run("topic creation with schema", func(t *testing.T) {
		fxtest.New(
			t,
			fx.NopLogger,
			fxconfig.FxConfigModule,
			fxgcppubsub.FxGcpPubSubModule,
			fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
			fxgcppubsub.PrepareTopicWithSchema(fxgcppubsub.PrepareTopicWithSchemaParams{
				TopicID:  "test-topic",
				SchemaID: "test-schema",
			}),
			fx.Populate(&factory),
		).RequireStart().RequireStop()

		sub, err := factory.Create(ctx, "test-topic")
		assert.NoError(t, err)

		assert.Equal(t, "test-topic", sub.BaseTopic().ID())
	})

	t.Run("topic creation error", func(t *testing.T) {
		fxtest.New(
			t,
			fx.NopLogger,
			fxconfig.FxConfigModule,
			fxgcppubsub.FxGcpPubSubModule,
			fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
			fx.Populate(&factory),
		).RequireStart().RequireStop()

		sub, err := factory.Create(ctx, "test-topic")
		assert.Nil(t, sub)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot get topic test-topic configuration")
	})
}
