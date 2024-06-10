package fxgcppubsub

import (
	"context"
	"fmt"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub/subscription"
)

type Subscriber struct {
	factory  *subscription.SubscriptionFactory
	registry *subscription.SubscriptionRegistry
}

func NewSubscriber(factory *subscription.SubscriptionFactory, registry *subscription.SubscriptionRegistry) *Subscriber {
	return &Subscriber{
		factory:  factory,
		registry: registry,
	}
}

func (s *Subscriber) Subscribe(ctx context.Context, subscriptionID string, f subscription.SubscribeFunc, options ...subscription.SubscribeOption) error {
	// get subscription
	if !s.registry.Has(subscriptionID) {
		subscription, err := s.factory.Create(ctx, subscriptionID)
		if err != nil {
			return err
		}

		s.registry.Add(subscription)
	}

	subscription, err := s.registry.Get(subscriptionID)
	if err != nil {
		return fmt.Errorf("cannot get subsciption: %w", err)
	}

	// subscribe
	return subscription.WithOptions(options...).Subscribe(ctx, f)
}
