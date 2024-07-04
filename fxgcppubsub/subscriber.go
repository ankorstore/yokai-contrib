package fxgcppubsub

import (
	"context"
	"fmt"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub/subscription"
)

var _ Subscriber = (*DefaultSubscriber)(nil)

// Subscriber is the interface for high level subscribers.
type Subscriber interface {
	Subscribe(ctx context.Context, subscriptionID string, f subscription.SubscribeFunc, options ...subscription.SubscribeOption) error
}

// DefaultSubscriber is the default Subscriber implementation.
type DefaultSubscriber struct {
	factory  subscription.SubscriptionFactory
	registry subscription.SubscriptionRegistry
}

// NewDefaultSubscriber returns a new DefaultSubscriber instance.
func NewDefaultSubscriber(factory subscription.SubscriptionFactory, registry subscription.SubscriptionRegistry) *DefaultSubscriber {
	return &DefaultSubscriber{
		factory:  factory,
		registry: registry,
	}
}

// Subscribe handle received data using a subscription.SubscribeFunc, with options, from a given subscriptionID.
func (s *DefaultSubscriber) Subscribe(ctx context.Context, subscriptionID string, f subscription.SubscribeFunc, options ...subscription.SubscribeOption) error {
	// retrieve subscription
	if !s.registry.Has(subscriptionID) {
		sub, err := s.factory.Create(ctx, subscriptionID)
		if err != nil {
			return fmt.Errorf("cannot create subscription: %w", err)
		}

		s.registry.Add(sub)
	}

	sub, err := s.registry.Get(subscriptionID)
	if err != nil {
		return fmt.Errorf("cannot get subscription: %w", err)
	}

	// subscribe
	return sub.WithOptions(options...).Subscribe(ctx, f)
}
