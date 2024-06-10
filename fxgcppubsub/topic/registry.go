package topic

import (
	"fmt"
	"sync"
)

type TopicRegistry struct {
	topics map[string]*Topic
	mutex  sync.RWMutex
}

func NewTopicRegistry() *TopicRegistry {
	return &TopicRegistry{
		topics: make(map[string]*Topic),
	}
}

func (r *TopicRegistry) Has(topicID string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, found := r.topics[topicID]

	return found
}

func (r *TopicRegistry) Add(topic *Topic) *TopicRegistry {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.topics[topic.Base().ID()] = topic

	return r
}

func (r *TopicRegistry) Get(topicID string) (*Topic, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, found := r.topics[topicID]; found {
		return r.topics[topicID], nil
	}

	return nil, fmt.Errorf("cannot find topic %s", topicID)
}

func (r *TopicRegistry) All() map[string]*Topic {
	return r.topics
}
