package fxgcppubsub_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/message"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/subscription"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type subscriptionFactoryMock struct {
	mock.Mock
}

func (m *subscriptionFactoryMock) Create(ctx context.Context, subscriptionID string) (*subscription.Subscription, error) {
	args := m.Called(ctx, subscriptionID)

	if t, ok := args.Get(0).(*subscription.Subscription); ok {
		return t, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

type subscriptionRegistryMock struct {
	mock.Mock
}

func (m *subscriptionRegistryMock) Has(subscriptionID string) bool {
	args := m.Called(subscriptionID)

	return args.Bool(0)
}

func (m *subscriptionRegistryMock) Get(subscriptionID string) (*subscription.Subscription, error) {
	args := m.Called(subscriptionID)

	if t, ok := args.Get(0).(*subscription.Subscription); ok {
		return t, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (m *subscriptionRegistryMock) Add(subscription *subscription.Subscription) {
	m.Called(subscription)
}

func (m *subscriptionRegistryMock) All() map[string]*subscription.Subscription {
	args := m.Called()

	if ts, ok := args.Get(0).(map[string]*subscription.Subscription); ok {
		return ts
	} else {
		return nil
	}
}

func TestSubscriber(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("subscription creation error", func(t *testing.T) {
		t.Parallel()

		sfm := new(subscriptionFactoryMock)
		sfm.On("Create", ctx, "test-subscription").Return(nil, assert.AnError).Once()

		srm := new(subscriptionRegistryMock)
		srm.On("Has", "test-subscription").Return(false).Once()

		subscriber := fxgcppubsub.NewDefaultSubscriber(sfm, srm)

		err := subscriber.Subscribe(ctx, "test-subscription", func(ctx context.Context, m *message.Message) {})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot create subscription")

		sfm.AssertExpectations(t)
		srm.AssertExpectations(t)
	})

	t.Run("subscription retrieval error", func(t *testing.T) {
		t.Parallel()

		sfm := new(subscriptionFactoryMock)
		sfm.AssertNotCalled(t, "Create")

		srm := new(subscriptionRegistryMock)
		srm.On("Has", "test-subscription").Return(true).Once()
		srm.AssertNotCalled(t, "Add")
		srm.On("Get", "test-subscription").Return(nil, assert.AnError).Once()

		subscriber := fxgcppubsub.NewDefaultSubscriber(sfm, srm)

		err := subscriber.Subscribe(ctx, "test-subscription", func(ctx context.Context, m *message.Message) {})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot get subscription")

		sfm.AssertExpectations(t)
		srm.AssertExpectations(t)
	})
}
