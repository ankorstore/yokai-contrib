package subscription_test

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/message"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/subscription"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/avro"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/proto"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestSubscription(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	var publisher fxgcppubsub.Publisher
	var client *pubsub.Client
	var supervisor *reactor.WaiterSupervisor

	ctx := context.Background()
	avroSchemaDefinition := getTestAvroSchemaDefinition(t)
	protoSchemaDefinition := getTestProtoSchemaDefinition(t)

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
		fx.Populate(&publisher, &client, &supervisor),
	).RequireStart().RequireStop()

	t.Run("getters", func(t *testing.T) {
		cod := codec.NewDefaultCodec(pubsub.SchemaTypeUnspecified, pubsub.EncodingUnspecified, "")
		baseSub := client.Subscription("raw-subscription")
		sub := subscription.NewSubscription(cod, baseSub)

		assert.Equal(t, cod, sub.Codec())
		assert.Equal(t, baseSub, sub.BaseSubscription())
	})

	t.Run("raw message", func(t *testing.T) {
		cod := codec.NewDefaultCodec(pubsub.SchemaTypeUnspecified, pubsub.EncodingUnspecified, "")
		baseSub := client.Subscription("raw-subscription")
		sub := subscription.NewSubscription(cod, baseSub)

		_, err := publisher.Publish(ctx, "raw-topic", []byte("raw data"))
		assert.NoError(t, err)

		waiter := supervisor.StartWaiter("projects/test-project/subscriptions/raw-subscription")

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
		cod := codec.NewDefaultCodec(pubsub.SchemaAvro, pubsub.EncodingBinary, avroSchemaDefinition)
		baseSub := client.Subscription("avro-subscription")
		sub := subscription.NewSubscription(cod, baseSub)

		_, err := publisher.Publish(ctx, "avro-topic", &avro.SimpleRecord{
			StringField:  "test avro",
			FloatField:   12.34,
			BooleanField: true,
		})
		assert.NoError(t, err)

		waiter := supervisor.StartWaiter("projects/test-project/subscriptions/avro-subscription")

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
		cod := codec.NewDefaultCodec(pubsub.SchemaProtocolBuffer, pubsub.EncodingBinary, protoSchemaDefinition)
		baseSub := client.Subscription("proto-subscription")
		sub := subscription.NewSubscription(cod, baseSub)

		_, err := publisher.Publish(ctx, "proto-topic", &proto.SimpleRecord{
			StringField:  "test proto",
			FloatField:   56.78,
			BooleanField: false,
		})
		assert.NoError(t, err)

		waiter := supervisor.StartWaiter("projects/test-project/subscriptions/proto-subscription")

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

func getTestAvroSchemaDefinition(tb testing.TB) string {
	tb.Helper()

	data, err := os.ReadFile("../testdata/avro/simple.avsc")
	assert.NoError(tb, err)

	return string(data)
}

func getTestProtoSchemaDefinition(tb testing.TB) string {
	tb.Helper()

	data, err := os.ReadFile("../testdata/proto/simple.proto")
	assert.NoError(tb, err)

	return string(data)
}
