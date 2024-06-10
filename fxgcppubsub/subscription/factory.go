package subscription

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/schema"
)

type SubscriptionFactory struct {
	client   *pubsub.Client
	registry *schema.SchemaRegistry
}

func NewSubscriptionFactory(client *pubsub.Client, registry *schema.SchemaRegistry) *SubscriptionFactory {
	return &SubscriptionFactory{
		client:   client,
		registry: registry,
	}
}

func (f *SubscriptionFactory) Create(ctx context.Context, subscriptionID string) (*Subscription, error) {
	// subscription
	subscription := f.client.Subscription(subscriptionID)

	// subscription config
	subscriptionConfig, err := subscription.Config(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get subscription configuration: %w", err)
	}

	// topic config
	topicConfig, err := subscriptionConfig.Topic.Config(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get subscription topic configuration: %w", err)
	}

	// topic schema settings
	topicSchemaSettings := topicConfig.SchemaSettings

	// topic schema config
	var topicSchemaConfig *pubsub.SchemaConfig
	if topicSchemaSettings != nil {
		// schema config
		topicSchemaConfig, err = f.registry.Get(ctx, topicSchemaSettings.Schema)
		if err != nil {
			return nil, err
		}
	}

	// codec
	topicCodec := codec.NewCodec(topicSchemaConfig, topicSchemaSettings)

	// create
	return NewSubscription(subscription, topicCodec), nil
}
