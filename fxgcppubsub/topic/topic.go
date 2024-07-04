package topic

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
)

// Topic represents a pub/sub topic with an associated codec.Codec.
type Topic struct {
	codec   codec.Codec
	topic   *pubsub.Topic
	options *Options
}

// NewTopic returns a new Topic instance.
func NewTopic(codec codec.Codec, topic *pubsub.Topic) *Topic {
	return &Topic{
		codec:   codec,
		topic:   topic,
		options: DefaultPublishOptions(),
	}
}

// Codec returns the topic associated codec.Codec.
func (t *Topic) Codec() codec.Codec {
	return t.codec
}

// BaseTopic returns the base pubsub.Topic.
func (t *Topic) BaseTopic() *pubsub.Topic {
	return t.topic
}

// WithOptions configures the topic with a list of PublishOption.
func (t *Topic) WithOptions(options ...PublishOption) *Topic {
	// resolve options
	for _, applyOpt := range options {
		applyOpt(t.options)
	}

	// apply options
	t.topic.PublishSettings = t.options.PublishSettings
	t.topic.EnableMessageOrdering = t.options.MessageSettings.OrderingKey != ""

	return t
}

// Publish publishes the provided data.
func (t *Topic) Publish(ctx context.Context, data any) (*pubsub.PublishResult, error) {
	// encode
	encodedData, err := t.codec.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("cannot encode data: %w", err)
	}

	// publish
	return t.topic.Publish(ctx, &pubsub.Message{
		Data:        encodedData,
		Attributes:  t.options.MessageSettings.Attributes,
		OrderingKey: t.options.MessageSettings.OrderingKey,
	}), nil
}
