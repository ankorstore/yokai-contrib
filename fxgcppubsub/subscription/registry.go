package subscription

import (
	"fmt"
	"sync"
)

var _ SubscriptionRegistry = (*DefaultSubscriptionRegistry)(nil)

// SubscriptionRegistry is the interface for Subscription registries.
type SubscriptionRegistry interface {
	Has(subscriptionID string) bool
	Get(subscriptionID string) (*Subscription, error)
	Add(subscription *Subscription)
	All() map[string]*Subscription
}

// DefaultSubscriptionRegistry is the default SubscriptionRegistry implementation.
type DefaultSubscriptionRegistry struct {
	subscriptions map[string]*Subscription
	mutex         sync.RWMutex
}

// NewDefaultSubscriptionRegistry returns a new DefaultSubscriptionRegistry instance.
func NewDefaultSubscriptionRegistry() *DefaultSubscriptionRegistry {
	return &DefaultSubscriptionRegistry{
		subscriptions: make(map[string]*Subscription),
	}
}

// Has returns true if the registry contains a Subscription for the provided subscriptionID.
func (r *DefaultSubscriptionRegistry) Has(subscriptionID string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, found := r.subscriptions[subscriptionID]

	return found
}

// Get returns a registered Subscription for the provided subscriptionID.
func (r *DefaultSubscriptionRegistry) Get(subscriptionID string) (*Subscription, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, found := r.subscriptions[subscriptionID]; found {
		return r.subscriptions[subscriptionID], nil
	}

	return nil, fmt.Errorf("cannot find subscription %s", subscriptionID)
}

// Add registers a Subscription.
func (r *DefaultSubscriptionRegistry) Add(subscription *Subscription) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.subscriptions[subscription.BaseSubscription().ID()] = subscription
}

// All returns all registered Subscription.
func (r *DefaultSubscriptionRegistry) All() map[string]*Subscription {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.subscriptions
}
