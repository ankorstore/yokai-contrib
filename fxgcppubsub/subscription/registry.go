package subscription

import (
	"fmt"
	"sync"
)

type SubscriptionRegistry struct {
	subscriptions map[string]*Subscription
	mutex         sync.RWMutex
}

func NewSubscriptionRegistry() *SubscriptionRegistry {
	return &SubscriptionRegistry{
		subscriptions: make(map[string]*Subscription),
	}
}

func (r *SubscriptionRegistry) Has(subscriptionID string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, found := r.subscriptions[subscriptionID]

	return found
}

func (r *SubscriptionRegistry) Add(subscription *Subscription) *SubscriptionRegistry {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.subscriptions[subscription.Base().ID()] = subscription

	return r
}

func (r *SubscriptionRegistry) Get(topicID string) (*Subscription, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, found := r.subscriptions[topicID]; found {
		return r.subscriptions[topicID], nil
	}

	return nil, fmt.Errorf("cannot find subscription %s", topicID)
}

func (r *SubscriptionRegistry) All() map[string]*Subscription {
	return r.subscriptions
}
