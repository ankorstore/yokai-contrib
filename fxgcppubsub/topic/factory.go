package topic

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/schema"
)

type TopicFactory struct {
	client   *pubsub.Client
	registry *schema.SchemaRegistry
}

func NewTopicFactory(client *pubsub.Client, registry *schema.SchemaRegistry) *TopicFactory {
	return &TopicFactory{
		client:   client,
		registry: registry,
	}
}

func (f *TopicFactory) Create(ctx context.Context, topicID string) (*Topic, error) {
	// topic
	topic := f.client.Topic(topicID)

	// topic config
	topicConfig, err := topic.Config(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get topic config: %w", err)
	}

	// topic schema settings
	topicSchemaSettings := topicConfig.SchemaSettings

	// topic schema config
	var topicSchemaConfig *pubsub.SchemaConfig
	if topicSchemaSettings != nil {
		topicSchemaConfig, err = f.registry.Get(ctx, topicSchemaSettings.Schema)
		if err != nil {
			return nil, err
		}
	}

	// codec
	topicCodec := codec.NewCodec(topicSchemaConfig, topicSchemaSettings)

	// create
	return NewTopic(topic, topicCodec), nil
}
