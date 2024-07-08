package topic_test

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/message"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor/ack"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/avro"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/proto"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/topic"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestTopic(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "../testdata/config")
	t.Setenv("GCP_PROJECT_ID", "test-project")

	var subscriber fxgcppubsub.Subscriber
	var supervisor ack.AckSupervisor
	var client *pubsub.Client

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
		fx.Populate(&subscriber, &client, &supervisor),
	).RequireStart().RequireStop()

	t.Run("getters", func(t *testing.T) {
		cod := codec.NewRawCodec()
		baseTop := client.Topic("raw-topic")
		top := topic.NewTopic(cod, baseTop)

		assert.Equal(t, cod, top.Codec())
		assert.Equal(t, baseTop, top.BaseTopic())
	})

	t.Run("raw message", func(t *testing.T) {
		cod := codec.NewRawCodec()
		baseTop := client.Topic("raw-topic")
		top := topic.NewTopic(cod, baseTop)

		res, err := top.Publish(ctx, []byte("raw data"))
		assert.NoError(t, err)

		sid, err := res.Get(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, sid)

		waiter := supervisor.StartAckWaiter("raw-subscription")

		var out []byte

		//nolint:errcheck
		go subscriber.Subscribe(ctx, "raw-subscription", func(ctx context.Context, m *message.Message) {
			out = m.Data()

			m.Ack()
		})

		_, err = waiter.WaitMaxDuration(ctx, 1*time.Second)
		assert.NoError(t, err)

		assert.Equal(t, []byte("raw data"), out)
	})

	t.Run("avro message", func(t *testing.T) {
		cod := codec.NewAvroBinaryCodec(avroSchemaDefinition)
		baseTop := client.Topic("avro-topic")
		top := topic.NewTopic(cod, baseTop)

		res, err := top.Publish(ctx, &avro.SimpleRecord{
			StringField:  "test avro",
			FloatField:   12.34,
			BooleanField: true,
		})
		assert.NoError(t, err)

		sid, err := res.Get(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, sid)

		waiter := supervisor.StartAckWaiter("avro-subscription")

		var out avro.SimpleRecord

		//nolint:errcheck
		go subscriber.Subscribe(ctx, "avro-subscription", func(ctx context.Context, m *message.Message) {
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
		cod := codec.NewProtoBinaryCodec(protoSchemaDefinition)
		baseTop := client.Topic("proto-topic")
		top := topic.NewTopic(cod, baseTop)

		res, err := top.Publish(ctx, &proto.SimpleRecord{
			StringField:  "test proto",
			FloatField:   56.78,
			BooleanField: false,
		})
		assert.NoError(t, err)

		sid, err := res.Get(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, sid)

		waiter := supervisor.StartAckWaiter("proto-subscription")

		var out proto.SimpleRecord

		//nolint:errcheck
		go subscriber.Subscribe(ctx, "proto-subscription", func(ctx context.Context, m *message.Message) {
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
