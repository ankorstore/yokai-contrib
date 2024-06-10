package topic

import (
	"time"

	"cloud.google.com/go/pubsub"
)

type Options struct {
	Settings pubsub.PublishSettings
}

func DefaultPublishOptions() Options {
	return Options{
		Settings: pubsub.DefaultPublishSettings,
	}
}

type PublishOption func(o *Options)

func WithDelayThreshold(t time.Duration) PublishOption {
	return func(o *Options) {
		o.Settings.DelayThreshold = t
	}
}

func WithCountThreshold(n int) PublishOption {
	return func(o *Options) {
		o.Settings.CountThreshold = n
	}
}

func WithByteThreshold(n int) PublishOption {
	return func(o *Options) {
		o.Settings.ByteThreshold = n
	}
}

func WithNumGoroutines(n int) PublishOption {
	return func(o *Options) {
		o.Settings.NumGoroutines = n
	}
}

func WithTimeout(t time.Duration) PublishOption {
	return func(o *Options) {
		o.Settings.Timeout = t
	}
}

func WithFlowControlSettings(s pubsub.FlowControlSettings) PublishOption {
	return func(o *Options) {
		o.Settings.FlowControlSettings = s
	}
}

func WithCompression(c bool) PublishOption {
	return func(o *Options) {
		o.Settings.EnableCompression = c
	}
}

func WithCompressionBytesThreshold(n int) PublishOption {
	return func(o *Options) {
		o.Settings.CompressionBytesThreshold = n
	}
}
