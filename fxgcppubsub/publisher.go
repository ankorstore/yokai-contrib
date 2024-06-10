package fxgcppubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/topic"
)

type Publisher struct {
	factory  *topic.TopicFactory
	registry *topic.TopicRegistry
}

func NewPublisher(factory *topic.TopicFactory, registry *topic.TopicRegistry) *Publisher {
	return &Publisher{
		factory:  factory,
		registry: registry,
	}
}

func (p *Publisher) Publish(ctx context.Context, topicID string, data any, options ...topic.PublishOption) (*pubsub.PublishResult, error) {
	// get topic
	if !p.registry.Has(topicID) {
		topic, err := p.factory.Create(ctx, topicID)
		if err != nil {
			return nil, err
		}

		p.registry.Add(topic)
	}

	topic, err := p.registry.Get(topicID)
	if err != nil {
		return nil, fmt.Errorf("cannot get topic: %w", err)
	}

	// publish
	return topic.WithOptions(options...).Publish(ctx, data)
}

func (p *Publisher) Stop() {
	// graceful stop
	for _, topic := range p.registry.All() {
		topic.Base().Stop()
	}
}
