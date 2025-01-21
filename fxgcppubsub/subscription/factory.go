package subscription

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/schema"
)

var _ SubscriptionFactory = (*DefaultSubscriptionFactory)(nil)

// SubscriptionFactory is the interface for Subscription factories.
type SubscriptionFactory interface {
	Create(ctx context.Context, subscriptionID string) (*Subscription, error)
}

// DefaultSubscriptionFactory is the default SubscriptionFactory implementation.
type DefaultSubscriptionFactory struct {
	client   *pubsub.Client
	registry schema.SchemaConfigRegistry
	factory  codec.CodecFactory
}

// NewDefaultSubscriptionFactory returns a new DefaultSubscriptionFactory instance.
func NewDefaultSubscriptionFactory(client *pubsub.Client, registry schema.SchemaConfigRegistry, factory codec.CodecFactory) *DefaultSubscriptionFactory {
	return &DefaultSubscriptionFactory{
		client:   client,
		registry: registry,
		factory:  factory,
	}
}

// Create creates a new Subscription.
func (f *DefaultSubscriptionFactory) Create(ctx context.Context, subscriptionID string) (*Subscription, error) {
	// subscription
	subscription := f.client.Subscription(subscriptionID)

	// subscription config
	subscriptionConfig, err := subscription.Config(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get subscription %s configuration: %w", subscriptionID, err)
	}

	// subscription topic config
	topicConfig, err := subscriptionConfig.Topic.Config(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get subscription %s topic configuration: %w", subscriptionID, err)
	}

	// subscription topic schema config
	topicSchemaType := pubsub.SchemaTypeUnspecified
	topicSchemaEncoding := pubsub.EncodingUnspecified
	topicSchemaDefinition := ""

	if topicConfig.SchemaSettings != nil {
		topicSchemaConfig, err := f.registry.Get(ctx, topicConfig.SchemaSettings.Schema)
		if err != nil {
			return nil, fmt.Errorf("cannot get subscription %s topic schema configuration: %w", subscriptionID, err)
		}

		topicSchemaType = topicSchemaConfig.Type
		topicSchemaEncoding = topicConfig.SchemaSettings.Encoding
		topicSchemaDefinition = topicSchemaConfig.Definition
	}

	// subscription codec
	subscriptionCodec, err := f.factory.Create(topicSchemaType, topicSchemaEncoding, topicSchemaDefinition)
	if err != nil {
		return nil, fmt.Errorf("cannot create subscription %s codec: %w", subscriptionID, err)
	}

	return NewSubscription(subscriptionCodec, subscription), nil
}
