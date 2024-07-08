package fxgcppubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"go.uber.org/fx"
)

// PrepareSchemaParams represents the parameters used in PrepareSchema.
type PrepareSchemaParams struct {
	SchemaID     string
	SchemaConfig pubsub.SchemaConfig
}

// PrepareSchema prepares a pub/sub schema.
func PrepareSchema(params PrepareSchemaParams) fx.Option {
	return fx.Invoke(
		func(ctx context.Context, schemaClient *pubsub.SchemaClient) error {
			_, err := schemaClient.CreateSchema(ctx, params.SchemaID, params.SchemaConfig)
			if err != nil {
				return fmt.Errorf("failed to prepare schema %q: %w", params.SchemaID, err)
			}

			return nil
		},
	)
}

// PrepareTopicParams represents the parameters used in PrepareTopic.
type PrepareTopicParams struct {
	TopicID     string
	TopicConfig pubsub.TopicConfig
}

// PrepareTopic prepares a pub/sub topic.
func PrepareTopic(params PrepareTopicParams) fx.Option {
	return fx.Invoke(
		func(ctx context.Context, client *pubsub.Client) error {
			_, err := client.CreateTopicWithConfig(ctx, params.TopicID, &params.TopicConfig)
			if err != nil {
				return fmt.Errorf("failed to prepare topic %q: %w", params.TopicID, err)
			}

			return nil
		},
	)
}

// PrepareTopicWithSchemaParams represents the parameters used in PrepareTopicWithSchema.
type PrepareTopicWithSchemaParams struct {
	TopicID        string
	TopicConfig    pubsub.TopicConfig
	SchemaID       string
	SchemaConfig   pubsub.SchemaConfig
	SchemaEncoding pubsub.SchemaEncoding
}

// PrepareTopicWithSchema prepares a pub/sub topic with a schema.
func PrepareTopicWithSchema(params PrepareTopicWithSchemaParams) fx.Option {
	return fx.Invoke(
		func(ctx context.Context, schemaClient *pubsub.SchemaClient, client *pubsub.Client) error {
			// prepare schema
			schema, err := schemaClient.CreateSchema(ctx, params.SchemaID, params.SchemaConfig)
			if err != nil {
				return fmt.Errorf("failed to prepare schema %q: %w", params.SchemaID, err)
			}

			// prepare topic
			topicConfig := params.TopicConfig
			topicConfig.SchemaSettings = &pubsub.SchemaSettings{
				Schema:   schema.Name,
				Encoding: params.SchemaEncoding,
			}

			_, err = client.CreateTopicWithConfig(ctx, params.TopicID, &topicConfig)
			if err != nil {
				return fmt.Errorf("failed to prepare topic %q: %w", params.TopicID, err)
			}

			return nil
		},
	)
}

// PrepareTopicAndSubscriptionParams represents the parameters used in PrepareTopicAndSubscription.
type PrepareTopicAndSubscriptionParams struct {
	TopicID            string
	TopicConfig        pubsub.TopicConfig
	SubscriptionID     string
	SubscriptionConfig pubsub.SubscriptionConfig
}

// PrepareTopicAndSubscription prepares a pub/sub topic and an associated subscription.
func PrepareTopicAndSubscription(params PrepareTopicAndSubscriptionParams) fx.Option {
	return fx.Invoke(
		func(ctx context.Context, client *pubsub.Client) error {
			// prepare topic
			topic, err := client.CreateTopicWithConfig(ctx, params.TopicID, &params.TopicConfig)
			if err != nil {
				return fmt.Errorf("failed to prepare topic %q: %w", params.TopicID, err)
			}

			// prepare subscription
			subscriptionConfig := params.SubscriptionConfig
			subscriptionConfig.Topic = topic

			_, err = client.CreateSubscription(ctx, params.SubscriptionID, subscriptionConfig)
			if err != nil {
				return fmt.Errorf("failed to prepare subscription %q: %w", params.SubscriptionID, err)
			}

			return nil
		},
	)
}

// PrepareTopicAndSubscriptionWithSchemaParams represents the parameters used in PrepareTopicAndSubscriptionWithSchema.
type PrepareTopicAndSubscriptionWithSchemaParams struct {
	TopicID            string
	TopicConfig        pubsub.TopicConfig
	SubscriptionID     string
	SubscriptionConfig pubsub.SubscriptionConfig
	SchemaID           string
	SchemaConfig       pubsub.SchemaConfig
	SchemaEncoding     pubsub.SchemaEncoding
}

// PrepareTopicAndSubscriptionWithSchema prepares a pub/sub topic and an associated subscription with a schema.
func PrepareTopicAndSubscriptionWithSchema(params PrepareTopicAndSubscriptionWithSchemaParams) fx.Option {
	return fx.Invoke(
		func(ctx context.Context, schemaClient *pubsub.SchemaClient, client *pubsub.Client) error {
			// prepare schema
			schema, err := schemaClient.CreateSchema(ctx, params.SchemaID, params.SchemaConfig)
			if err != nil {
				return fmt.Errorf("failed to prepare schema %q: %w", params.SchemaID, err)
			}

			// prepare topic
			topicConfig := params.TopicConfig
			topicConfig.SchemaSettings = &pubsub.SchemaSettings{
				Schema:   schema.Name,
				Encoding: params.SchemaEncoding,
			}

			topic, err := client.CreateTopicWithConfig(ctx, params.TopicID, &topicConfig)
			if err != nil {
				return fmt.Errorf("failed to prepare topic %q: %w", params.TopicID, err)
			}

			// prepare subscription
			subscriptionConfig := params.SubscriptionConfig
			subscriptionConfig.Topic = topic

			_, err = client.CreateSubscription(ctx, params.SubscriptionID, subscriptionConfig)
			if err != nil {
				return fmt.Errorf("failed to prepare subscription %q: %w", params.SubscriptionID, err)
			}

			return nil
		},
	)
}
