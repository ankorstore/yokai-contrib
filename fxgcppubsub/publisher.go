package fxgcppubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/topic"
)

var _ Publisher = (*DefaultPublisher)(nil)

// Publisher is the interface for high level publishers.
type Publisher interface {
	Publish(ctx context.Context, topicID string, data any, options ...topic.PublishOption) (*pubsub.PublishResult, error)
	Stop()
}

// DefaultPublisher is the default Publisher implementation.
type DefaultPublisher struct {
	factory  topic.TopicFactory
	registry topic.TopicRegistry
}

// NewDefaultPublisher returns a new DefaultPublisher instance.
func NewDefaultPublisher(factory topic.TopicFactory, registry topic.TopicRegistry) *DefaultPublisher {
	return &DefaultPublisher{
		factory:  factory,
		registry: registry,
	}
}

// Publish publishes data, with options, on a given topicID.
func (p *DefaultPublisher) Publish(ctx context.Context, topicID string, data any, options ...topic.PublishOption) (*pubsub.PublishResult, error) {
	// retrieve topic
	if !p.registry.Has(topicID) {
		top, err := p.factory.Create(ctx, topicID)
		if err != nil {
			return nil, fmt.Errorf("cannot create topic: %w", err)
		}

		p.registry.Add(top)
	}

	top, err := p.registry.Get(topicID)
	if err != nil {
		return nil, fmt.Errorf("cannot get topic: %w", err)
	}

	// publish
	return top.WithOptions(options...).Publish(ctx, data)
}

// Stop stops gracefully all internal publishers.
func (p *DefaultPublisher) Stop() {
	for _, top := range p.registry.All() {
		top.BaseTopic().Stop()
	}
}
