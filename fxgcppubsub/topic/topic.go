package topic

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
)

type Topic struct {
	topic *pubsub.Topic
	codec *codec.Codec
}

func NewTopic(topic *pubsub.Topic, codec *codec.Codec) *Topic {
	return &Topic{
		topic: topic,
		codec: codec,
	}
}

func (t *Topic) Base() *pubsub.Topic {
	return t.topic
}

func (t *Topic) WithOptions(options ...PublishOption) *Topic {
	// resolve options
	topicOptions := DefaultPublishOptions()
	for _, applyOpt := range options {
		applyOpt(&topicOptions)
	}

	// apply options
	t.topic.PublishSettings = topicOptions.Settings

	return t
}

func (t *Topic) Publish(ctx context.Context, data any) (*pubsub.PublishResult, error) {
	// encode
	encodedData, err := t.codec.Encode(data)
	if err != nil {
		return nil, err
	}

	// publish
	return t.topic.Publish(ctx, &pubsub.Message{
		Data: encodedData,
	}), nil
}
