package topic_test

import (
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/topic"
	"github.com/stretchr/testify/assert"
)

func TestPublishOptions(t *testing.T) {
	t.Parallel()

	t.Run("withDelayThreshold", func(t *testing.T) {
		t.Parallel()

		o := &topic.Options{}
		value := time.Duration(1)
		opt := topic.WithDelayThreshold(value)
		opt(o)

		assert.Equal(t, value, o.PublishSettings.DelayThreshold)
	})

	t.Run("withCountThreshold", func(t *testing.T) {
		t.Parallel()

		o := &topic.Options{}
		value := 1
		opt := topic.WithCountThreshold(value)
		opt(o)

		assert.Equal(t, value, o.PublishSettings.CountThreshold)
	})

	t.Run("withByteThreshold", func(t *testing.T) {
		t.Parallel()

		o := &topic.Options{}
		value := 2
		opt := topic.WithByteThreshold(value)
		opt(o)

		assert.Equal(t, value, o.PublishSettings.ByteThreshold)
	})

	t.Run("withNumGoroutines", func(t *testing.T) {
		t.Parallel()

		o := &topic.Options{}
		value := 3
		opt := topic.WithNumGoroutines(value)
		opt(o)

		assert.Equal(t, value, o.PublishSettings.NumGoroutines)
	})

	t.Run("withTimeout", func(t *testing.T) {
		t.Parallel()

		o := &topic.Options{}
		value := time.Duration(2)
		opt := topic.WithTimeout(value)
		opt(o)

		assert.Equal(t, value, o.PublishSettings.Timeout)
	})

	t.Run("withFlowControlSettings", func(t *testing.T) {
		t.Parallel()

		o := &topic.Options{}
		value := pubsub.FlowControlSettings{}
		opt := topic.WithFlowControlSettings(value)
		opt(o)

		assert.Equal(t, value, o.PublishSettings.FlowControlSettings)
	})

	t.Run("WithCompression", func(t *testing.T) {
		t.Parallel()

		o := &topic.Options{}
		opt := topic.WithCompression(true)
		opt(o)

		assert.True(t, o.PublishSettings.EnableCompression)
	})

	t.Run("withCompressionBytesThreshold", func(t *testing.T) {
		t.Parallel()

		o := &topic.Options{}
		value := 4
		opt := topic.WithCompressionBytesThreshold(value)
		opt(o)

		assert.Equal(t, value, o.PublishSettings.CompressionBytesThreshold)
	})

	t.Run("withMessageOrderingKey", func(t *testing.T) {
		t.Parallel()

		o := &topic.Options{}
		value := "test"
		opt := topic.WithMessageOrderingKey(value)
		opt(o)

		assert.Equal(t, value, o.MessageSettings.OrderingKey)
	})

	t.Run("withMessageAttributes", func(t *testing.T) {
		t.Parallel()

		o := &topic.Options{}
		value := map[string]string{"test-key": "test-value"}
		opt := topic.WithMessageAttributes(value)
		opt(o)

		assert.Equal(t, value, o.MessageSettings.Attributes)
	})
}
