package topic

import (
	"time"

	"cloud.google.com/go/pubsub"
)

// MessageSettings represents message publish options.
type MessageSettings struct {
	OrderingKey string
	Attributes  map[string]string
}

// Options represents publish options.
type Options struct {
	PublishSettings pubsub.PublishSettings
	MessageSettings MessageSettings
}

// DefaultPublishOptions is the default publish options.
func DefaultPublishOptions() *Options {
	return &Options{
		PublishSettings: pubsub.DefaultPublishSettings,
		MessageSettings: MessageSettings{
			OrderingKey: "",
			Attributes:  make(map[string]string),
		},
	}
}

// PublishOption represents publish functional options.
type PublishOption func(o *Options)

// WithDelayThreshold sets the delay threshold.
func WithDelayThreshold(t time.Duration) PublishOption {
	return func(o *Options) {
		o.PublishSettings.DelayThreshold = t
	}
}

// WithCountThreshold sets the count threshold.
func WithCountThreshold(n int) PublishOption {
	return func(o *Options) {
		o.PublishSettings.CountThreshold = n
	}
}

// WithByteThreshold sets the byte threshold.
func WithByteThreshold(n int) PublishOption {
	return func(o *Options) {
		o.PublishSettings.ByteThreshold = n
	}
}

// WithNumGoroutines sets the num of goroutines.
func WithNumGoroutines(n int) PublishOption {
	return func(o *Options) {
		o.PublishSettings.NumGoroutines = n
	}
}

// WithTimeout sets the timeout.
func WithTimeout(t time.Duration) PublishOption {
	return func(o *Options) {
		o.PublishSettings.Timeout = t
	}
}

// WithFlowControlSettings sets the flow control settings.
func WithFlowControlSettings(s pubsub.FlowControlSettings) PublishOption {
	return func(o *Options) {
		o.PublishSettings.FlowControlSettings = s
	}
}

// WithCompression sets the compression usage.
func WithCompression(c bool) PublishOption {
	return func(o *Options) {
		o.PublishSettings.EnableCompression = c
	}
}

// WithCompressionBytesThreshold sets the compression bytes threshold.
func WithCompressionBytesThreshold(n int) PublishOption {
	return func(o *Options) {
		o.PublishSettings.CompressionBytesThreshold = n
	}
}

// WithMessageOrderingKey sets the message ordering key.
func WithMessageOrderingKey(k string) PublishOption {
	return func(o *Options) {
		o.MessageSettings.OrderingKey = k
	}
}

// WithMessageAttributes sets the message attributes.
func WithMessageAttributes(a map[string]string) PublishOption {
	return func(o *Options) {
		o.MessageSettings.Attributes = a
	}
}
