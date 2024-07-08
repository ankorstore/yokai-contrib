package fxgcppubsub_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/topic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type topicFactoryMock struct {
	mock.Mock
}

func (m *topicFactoryMock) Create(ctx context.Context, topicID string) (*topic.Topic, error) {
	args := m.Called(ctx, topicID)

	if t, ok := args.Get(0).(*topic.Topic); ok {
		return t, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

type topicRegistryMock struct {
	mock.Mock
}

func (m *topicRegistryMock) Has(topicID string) bool {
	args := m.Called(topicID)

	return args.Bool(0)
}

func (m *topicRegistryMock) Get(topicID string) (*topic.Topic, error) {
	args := m.Called(topicID)

	if t, ok := args.Get(0).(*topic.Topic); ok {
		return t, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (m *topicRegistryMock) Add(topic *topic.Topic) {
	m.Called(topic)
}

func (m *topicRegistryMock) All() map[string]*topic.Topic {
	args := m.Called()

	if ts, ok := args.Get(0).(map[string]*topic.Topic); ok {
		return ts
	} else {
		return nil
	}
}

func TestPublisher(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("topic creation error", func(t *testing.T) {
		t.Parallel()

		tfm := new(topicFactoryMock)
		tfm.On("Create", ctx, "test-topic").Return(nil, assert.AnError).Once()

		trm := new(topicRegistryMock)
		trm.On("Has", "test-topic").Return(false).Once()

		publisher := fxgcppubsub.NewDefaultPublisher(tfm, trm)

		res, err := publisher.Publish(ctx, "test-topic", []byte("test"))
		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot create topic")

		tfm.AssertExpectations(t)
		trm.AssertExpectations(t)
	})

	t.Run("topic retrieval error", func(t *testing.T) {
		t.Parallel()

		tfm := new(topicFactoryMock)
		tfm.AssertNotCalled(t, "Create")

		trm := new(topicRegistryMock)
		trm.On("Has", "test-topic").Return(true).Once()
		trm.AssertNotCalled(t, "Add")
		trm.On("Get", "test-topic").Return(nil, assert.AnError).Once()

		publisher := fxgcppubsub.NewDefaultPublisher(tfm, trm)

		res, err := publisher.Publish(ctx, "test-topic", []byte("test"))
		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot get topic")

		tfm.AssertExpectations(t)
		trm.AssertExpectations(t)
	})
}
