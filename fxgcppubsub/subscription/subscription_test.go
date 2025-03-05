package subscription_test

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/message"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor/ack"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/subscription"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/avro"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/proto"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestSubscription(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	var publisher fxgcppubsub.Publisher
	var supervisor ack.AckSupervisor
	var client *pubsub.Client

	ctx := context.Background()
	avroSchemaDefinition := avro.GetTestAvroSchemaDefinition(t)
	protoSchemaDefinition := proto.GetTestProtoSchemaDefinition(t)

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
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
		fx.Populate(&publisher, &client, &supervisor),
	).RequireStart().RequireStop()

	t.Run("getters", func(t *testing.T) {
		cod := codec.NewRawCodec()
		baseSub := client.Subscription("raw-subscription")
		sub := subscription.NewSubscription(cod, baseSub)

		assert.Equal(t, cod, sub.Codec())
		assert.Equal(t, baseSub, sub.BaseSubscription())
	})

	t.Run("raw message", func(t *testing.T) {
		cod := codec.NewRawCodec()
		baseSub := client.Subscription("raw-subscription")
		sub := subscription.NewSubscription(cod, baseSub)

		_, err := publisher.Publish(ctx, "raw-topic", []byte("raw data"))
		assert.NoError(t, err)

		waiter := supervisor.StartAckWaiter("raw-subscription")

		var out []byte

		//nolint:errcheck
		go sub.
			WithOptions(
				subscription.WithNumGoroutines(1),
				subscription.WithMaxOutstandingMessages(1),
			).
			Subscribe(ctx, func(ctx context.Context, m *message.Message) {
				out = m.Data()

				m.Ack()
			})

		_, err = waiter.WaitMaxDuration(ctx, 1*time.Second)
		assert.NoError(t, err)

		assert.Equal(t, []byte("raw data"), out)
	})

	t.Run("avro message", func(t *testing.T) {
		cod, err := codec.NewAvroBinaryCodec(avroSchemaDefinition)
		assert.NoError(t, err)

		baseSub := client.Subscription("avro-subscription")
		sub := subscription.NewSubscription(cod, baseSub)

		_, err = publisher.Publish(ctx, "avro-topic", &avro.SimpleRecord{
			StringField:  "test avro",
			FloatField:   12.34,
			BooleanField: true,
		})
		assert.NoError(t, err)

		waiter := supervisor.StartAckWaiter("avro-subscription")

		var out avro.SimpleRecord

		//nolint:errcheck
		go sub.
			WithOptions(
				subscription.WithNumGoroutines(1),
				subscription.WithMaxOutstandingMessages(1),
			).
			Subscribe(ctx, func(ctx context.Context, m *message.Message) {
				err = m.Decode(&out)
				assert.NoError(t, err)

				m.Ack()
			})

		_, err = waiter.WaitMaxDuration(ctx, 1*time.Second)
		assert.NoError(t, err)

		assert.Equal(t, "test avro", out.StringField)
		assert.Equal(t, float32(12.34), out.FloatField)
		assert.True(t, out.BooleanField)
	})

	t.Run("proto message", func(t *testing.T) {
		cod := codec.NewProtoBinaryCodec()
		baseSub := client.Subscription("proto-subscription")
		sub := subscription.NewSubscription(cod, baseSub)

		_, err := publisher.Publish(ctx, "proto-topic", &proto.SimpleRecord{
			StringField:  "test proto",
			FloatField:   56.78,
			BooleanField: false,
		})
		assert.NoError(t, err)

		waiter := supervisor.StartAckWaiter("proto-subscription")

		var out proto.SimpleRecord

		//nolint:errcheck
		go sub.
			WithOptions(
				subscription.WithNumGoroutines(1),
				subscription.WithMaxOutstandingMessages(1),
			).
			Subscribe(ctx, func(ctx context.Context, m *message.Message) {
				err = m.Decode(&out)
				assert.NoError(t, err)

				m.Ack()
			})

		_, err = waiter.WaitMaxDuration(ctx, 1*time.Second)
		assert.NoError(t, err)

		assert.Equal(t, "test proto", out.StringField)
		assert.Equal(t, float32(56.78), out.FloatField)
		assert.False(t, out.BooleanField)
	})
}
