package subscription_test

import (
	"testing"
	"time"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub/subscription"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeOptions(t *testing.T) {
	t.Parallel()

	t.Run("withMaxExtension", func(t *testing.T) {
		t.Parallel()

		o := &subscription.Options{}
		value := time.Duration(1)
		opt := subscription.WithMaxExtension(value)
		opt(o)

		assert.Equal(t, value, o.ReceiveSettings.MaxExtension)
	})

	t.Run("withMinExtensionPeriod", func(t *testing.T) {
		t.Parallel()

		o := &subscription.Options{}
		value := time.Duration(2)
		opt := subscription.WithMinExtensionPeriod(value)
		opt(o)

		assert.Equal(t, value, o.ReceiveSettings.MinExtensionPeriod)
	})

	t.Run("withMaxExtensionPeriod", func(t *testing.T) {
		t.Parallel()

		o := &subscription.Options{}
		value := time.Duration(3)
		opt := subscription.WithMaxExtensionPeriod(value)
		opt(o)

		assert.Equal(t, value, o.ReceiveSettings.MaxExtensionPeriod)
	})

	t.Run("withMaxOutstandingMessages", func(t *testing.T) {
		t.Parallel()

		o := &subscription.Options{}
		value := 1
		opt := subscription.WithMaxOutstandingMessages(value)
		opt(o)

		assert.Equal(t, value, o.ReceiveSettings.MaxOutstandingMessages)
	})

	t.Run("withMaxOutstandingBytes", func(t *testing.T) {
		t.Parallel()

		o := &subscription.Options{}
		value := 2
		opt := subscription.WithMaxOutstandingBytes(value)
		opt(o)

		assert.Equal(t, value, o.ReceiveSettings.MaxOutstandingBytes)
	})

	t.Run("withMaxOutstandingBytes", func(t *testing.T) {
		t.Parallel()

		o := &subscription.Options{}
		opt := subscription.WithLegacyFlowControl(true)
		opt(o)

		assert.True(t, o.ReceiveSettings.UseLegacyFlowControl)
	})

	t.Run("WithNumGoroutines", func(t *testing.T) {
		t.Parallel()

		o := &subscription.Options{}
		value := 3
		opt := subscription.WithNumGoroutines(value)
		opt(o)

		assert.Equal(t, value, o.ReceiveSettings.NumGoroutines)
	})
}
