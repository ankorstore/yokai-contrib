package topic

import (
	"fmt"
	"sync"
)

var _ TopicRegistry = (*DefaultTopicRegistry)(nil)

// TopicRegistry is the interface for Topic registries.
type TopicRegistry interface {
	Has(topicID string) bool
	Get(topicID string) (*Topic, error)
	Add(topic *Topic)
	All() map[string]*Topic
}

// DefaultTopicRegistry is the default TopicRegistry implementation.
type DefaultTopicRegistry struct {
	topics map[string]*Topic
	mutex  sync.RWMutex
}

// NewDefaultTopicRegistry returns a new DefaultTopicRegistry instance.
func NewDefaultTopicRegistry() *DefaultTopicRegistry {
	return &DefaultTopicRegistry{
		topics: make(map[string]*Topic),
	}
}

// Has returns true if the registry contains a Topic for the provided topicID.
func (r *DefaultTopicRegistry) Has(topicID string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, found := r.topics[topicID]

	return found
}

// Get returns a registered Topic for the provided topicID.
func (r *DefaultTopicRegistry) Get(topicID string) (*Topic, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, found := r.topics[topicID]; found {
		return r.topics[topicID], nil
	}

	return nil, fmt.Errorf("cannot find topic %s", topicID)
}

// Add registers a Topic.
func (r *DefaultTopicRegistry) Add(topic *Topic) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.topics[topic.BaseTopic().ID()] = topic
}

// All returns all registered Topic.
func (r *DefaultTopicRegistry) All() map[string]*Topic {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.topics
}
