package fxgcppubsub_test

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/message"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/avro"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/proto"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestFxGcpPubSubModule(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	var publisher fxgcppubsub.Publisher
	var subscriber fxgcppubsub.Subscriber
	var supervisor *reactor.WaiterSupervisor

	ctx := context.Background()
	avroSchemaDefinition := avro.GetTestAvroSchemaDefinition(t)
	protoSchemaDefinition := proto.GetTestProtoSchemaDefinition(t)

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxgcppubsub.FxGcpPubSubModule,
		fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),
		fxgcppubsub.PrepareTopicAndSubscription(fxgcppubsub.PrepareTopicAndSubscriptionParams{
			TopicID:        "raw-topic",
			SubscriptionID: "raw-subscription",
		}),
		fxgcppubsub.PrepareTopicAndSubscriptionWithSchema(fxgcppubsub.PrepareTopicAndSubscriptionWithSchemaParams{
			TopicID:        "avro-topic",
			SubscriptionID: "avro-subscription",
			SchemaID:       "avro-schema",
			SchemaConfig: pubsub.SchemaConfig{
				Name:       "avro-schema",
				Type:       pubsub.SchemaAvro,
				Definition: avroSchemaDefinition,
			},
			SchemaEncoding: pubsub.EncodingBinary,
		}),
		fxgcppubsub.PrepareTopicAndSubscriptionWithSchema(fxgcppubsub.PrepareTopicAndSubscriptionWithSchemaParams{
			TopicID:        "proto-topic",
			SubscriptionID: "proto-subscription",
			SchemaID:       "proto-schema",
			SchemaConfig: pubsub.SchemaConfig{
				Name:       "proto-schema",
				Type:       pubsub.SchemaProtocolBuffer,
				Definition: protoSchemaDefinition,
			},
			SchemaEncoding: pubsub.EncodingBinary,
		}),
		fx.Populate(&publisher, &subscriber, &supervisor),
	).RequireStart().RequireStop()

	t.Run("raw message", func(t *testing.T) {
		res, err := publisher.Publish(ctx, "raw-topic", []byte("test"))
		assert.NotNil(t, res)
		assert.NoError(t, err)

		sid, err := res.Get(ctx)
		assert.NotEmpty(t, sid)
		assert.NoError(t, err)

		waiter := supervisor.StartWaiter("projects/test-project/subscriptions/raw-subscription")

		//nolint:errcheck
		go subscriber.Subscribe(ctx, "raw-subscription", func(ctx context.Context, m *message.Message) {
			assert.Equal(t, []byte("test"), m.Data())

			m.Ack()
		})

		_, err = waiter.WaitMaxDuration(ctx, time.Second)
		assert.NoError(t, err)
	})

	t.Run("avro message", func(t *testing.T) {
		res, err := publisher.Publish(ctx, "avro-topic", &avro.SimpleRecord{
			StringField:  "test avro",
			FloatField:   12.34,
			BooleanField: true,
		})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		sid, err := res.Get(ctx)
		assert.NotEmpty(t, sid)
		assert.NoError(t, err)

		waiter := supervisor.StartWaiter("projects/test-project/subscriptions/avro-subscription")

		//nolint:errcheck
		go subscriber.Subscribe(ctx, "avro-subscription", func(ctx context.Context, m *message.Message) {
			var out avro.SimpleRecord

			err = m.Decode(&out)
			assert.NoError(t, err)

			assert.Equal(t, "test avro", out.StringField)
			assert.Equal(t, float32(12.34), out.FloatField)
			assert.True(t, out.BooleanField)

			m.Ack()
		})

		_, err = waiter.WaitMaxDuration(ctx, time.Second)
		assert.NoError(t, err)
	})

	t.Run("proto message", func(t *testing.T) {
		res, err := publisher.Publish(ctx, "proto-topic", &proto.SimpleRecord{
			StringField:  "test proto",
			FloatField:   56.78,
			BooleanField: false,
		})
		assert.NotNil(t, res)
		assert.NoError(t, err)

		sid, err := res.Get(ctx)
		assert.NotEmpty(t, sid)
		assert.NoError(t, err)

		waiter := supervisor.StartWaiter("projects/test-project/subscriptions/proto-subscription")

		//nolint:errcheck
		go subscriber.Subscribe(ctx, "proto-subscription", func(ctx context.Context, m *message.Message) {
			var out proto.SimpleRecord

			err = m.Decode(&out)
			assert.NoError(t, err)

			assert.Equal(t, "test proto", out.StringField)
			assert.Equal(t, float32(56.78), out.FloatField)
			assert.False(t, out.BooleanField)

			m.Ack()
		})

		_, err = waiter.WaitMaxDuration(ctx, time.Second)
		assert.NoError(t, err)
	})
}
