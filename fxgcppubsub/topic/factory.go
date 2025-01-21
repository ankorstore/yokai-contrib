package topic

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/schema"
)

var _ TopicFactory = (*DefaultTopicFactory)(nil)

// TopicFactory is the interface for Topic factories.
type TopicFactory interface {
	Create(ctx context.Context, topicID string) (*Topic, error)
}

// DefaultTopicFactory is the default TopicFactory implementation.
type DefaultTopicFactory struct {
	client   *pubsub.Client
	registry schema.SchemaConfigRegistry
	factory  codec.CodecFactory
}

// NewDefaultTopicFactory returns a new DefaultTopicFactory instance.
func NewDefaultTopicFactory(client *pubsub.Client, registry schema.SchemaConfigRegistry, factory codec.CodecFactory) *DefaultTopicFactory {
	return &DefaultTopicFactory{
		client:   client,
		registry: registry,
		factory:  factory,
	}
}

// Create creates a new Topic.
func (f *DefaultTopicFactory) Create(ctx context.Context, topicID string) (*Topic, error) {
	// topic
	topic := f.client.Topic(topicID)

	// topic config
	topicConfig, err := topic.Config(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get topic %s configuration: %w", topicID, err)
	}

	// topic schema config
	topicSchemaType := pubsub.SchemaTypeUnspecified
	topicSchemaEncoding := pubsub.EncodingUnspecified
	topicSchemaDefinition := ""

	if topicConfig.SchemaSettings != nil {
		topicSchemaConfig, err := f.registry.Get(ctx, topicConfig.SchemaSettings.Schema)
		if err != nil {
			return nil, fmt.Errorf("cannot get topic %s schema configuration: %w", topicID, err)
		}

		topicSchemaType = topicSchemaConfig.Type
		topicSchemaEncoding = topicConfig.SchemaSettings.Encoding
		topicSchemaDefinition = topicSchemaConfig.Definition
	}

	// topic codec
	topicCodec, err := f.factory.Create(topicSchemaType, topicSchemaEncoding, topicSchemaDefinition)
	if err != nil {
		return nil, fmt.Errorf("cannot create topic %s codec: %w", topicID, err)
	}

	return NewTopic(topicCodec, topic), nil
}
