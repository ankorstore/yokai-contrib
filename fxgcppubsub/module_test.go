package fxgcppubsub_test

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/message"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor/ack"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/avro"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/proto"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestFxGcpPubSubModule(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	runTest := func(tb testing.TB) (context.Context, fxgcppubsub.Publisher, fxgcppubsub.Subscriber, ack.AckSupervisor) {
		tb.Helper()

		var publisher fxgcppubsub.Publisher
		var subscriber fxgcppubsub.Subscriber
		var supervisor ack.AckSupervisor

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
			fx.Populate(&publisher, &subscriber, &supervisor),
		).RequireStart().RequireStop()

		return ctx, publisher, subscriber, supervisor
	}

	t.Run("raw message ack", func(t *testing.T) {
		ctx, publisher, subscriber, supervisor := runTest(t)

		_, err := publisher.Publish(ctx, "raw-topic", []byte("test"))
		assert.NoError(t, err)

		publisher.Stop()

		waiter := supervisor.StartAckWaiter("raw-subscription")

		//nolint:errcheck
		go subscriber.Subscribe(ctx, "raw-subscription", func(ctx context.Context, m *message.Message) {
			assert.Equal(t, []byte("test"), m.Data())

			m.Ack()
		})

		_, err = waiter.WaitMaxDuration(ctx, 2*time.Second)
		assert.NoError(t, err)
	})

	t.Run("raw message nack", func(t *testing.T) {
		ctx, publisher, subscriber, supervisor := runTest(t)

		_, err := publisher.Publish(ctx, "raw-topic", []byte("test"))
		assert.NoError(t, err)

		publisher.Stop()

		waiter := supervisor.StartNackWaiter("raw-subscription")

		//nolint:errcheck
		go subscriber.Subscribe(ctx, "raw-subscription", func(ctx context.Context, m *message.Message) {
			assert.Equal(t, []byte("test"), m.Data())

			m.Nack()
		})

		_, err = waiter.WaitMaxDuration(ctx, 2*time.Second)
		assert.NoError(t, err)
	})

	t.Run("avro message ack", func(t *testing.T) {
		ctx, publisher, subscriber, supervisor := runTest(t)

		_, err := publisher.Publish(ctx, "avro-topic", &avro.SimpleRecord{
			StringField:  "test avro",
			FloatField:   12.34,
			BooleanField: true,
		})
		assert.NoError(t, err)

		publisher.Stop()

		waiter := supervisor.StartAckWaiter("avro-subscription")

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

		_, err = waiter.WaitMaxDuration(ctx, 2*time.Second)
		assert.NoError(t, err)
	})

	t.Run("avro message nack", func(t *testing.T) {
		ctx, publisher, subscriber, supervisor := runTest(t)

		_, err := publisher.Publish(ctx, "avro-topic", &avro.SimpleRecord{
			StringField:  "test avro",
			FloatField:   12.34,
			BooleanField: true,
		})
		assert.NoError(t, err)

		publisher.Stop()

		waiter := supervisor.StartNackWaiter("avro-subscription")

		//nolint:errcheck
		go subscriber.Subscribe(ctx, "avro-subscription", func(ctx context.Context, m *message.Message) {
			var out avro.SimpleRecord

			err = m.Decode(&out)
			assert.NoError(t, err)

			assert.Equal(t, "test avro", out.StringField)
			assert.Equal(t, float32(12.34), out.FloatField)
			assert.True(t, out.BooleanField)

			m.Nack()
		})

		_, err = waiter.WaitMaxDuration(ctx, 2*time.Second)
		assert.NoError(t, err)
	})

	t.Run("proto message ack", func(t *testing.T) {
		ctx, publisher, subscriber, supervisor := runTest(t)

		_, err := publisher.Publish(ctx, "proto-topic", &proto.SimpleRecord{
			StringField:  "test proto",
			FloatField:   56.78,
			BooleanField: false,
		})
		assert.NoError(t, err)

		publisher.Stop()

		waiter := supervisor.StartAckWaiter("proto-subscription")

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

		_, err = waiter.WaitMaxDuration(ctx, 2*time.Second)
		assert.NoError(t, err)
	})

	t.Run("proto message nack", func(t *testing.T) {
		ctx, publisher, subscriber, supervisor := runTest(t)

		_, err := publisher.Publish(ctx, "proto-topic", &proto.SimpleRecord{
			StringField:  "test proto",
			FloatField:   56.78,
			BooleanField: false,
		})
		assert.NoError(t, err)

		publisher.Stop()

		waiter := supervisor.StartNackWaiter("proto-subscription")

		//nolint:errcheck
		go subscriber.Subscribe(ctx, "proto-subscription", func(ctx context.Context, m *message.Message) {
			var out proto.SimpleRecord

			err = m.Decode(&out)
			assert.NoError(t, err)

			assert.Equal(t, "test proto", out.StringField)
			assert.Equal(t, float32(56.78), out.FloatField)
			assert.False(t, out.BooleanField)

			m.Nack()
		})

		_, err = waiter.WaitMaxDuration(ctx, 2*time.Second)
		assert.NoError(t, err)
	})
}
