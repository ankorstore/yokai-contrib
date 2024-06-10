package subscription

import (
	"time"

	"cloud.google.com/go/pubsub"
)

type Options struct {
	Settings pubsub.ReceiveSettings
}

func DefaultSubscribeOptions() Options {
	return Options{
		Settings: pubsub.DefaultReceiveSettings,
	}

}

type SubscribeOption func(o *Options)

func WithMaxExtension(t time.Duration) SubscribeOption {
	return func(o *Options) {
		o.Settings.MaxExtension = t
	}
}

func WithMinExtensionPeriod(t time.Duration) SubscribeOption {
	return func(o *Options) {
		o.Settings.MinExtensionPeriod = t
	}
}

func WithMaxExtensionPeriod(t time.Duration) SubscribeOption {
	return func(o *Options) {
		o.Settings.MaxExtensionPeriod = t
	}
}

func WithMaxOutstandingMessages(n int) SubscribeOption {
	return func(o *Options) {
		o.Settings.MaxOutstandingMessages = n
	}
}

func WithMaxOutstandingBytes(n int) SubscribeOption {
	return func(o *Options) {
		o.Settings.MaxOutstandingBytes = n
	}
}

func WithLegacyFlowControl(c bool) SubscribeOption {
	return func(o *Options) {
		o.Settings.UseLegacyFlowControl = c
	}
}

func WithNumGoroutines(n int) SubscribeOption {
	return func(o *Options) {
		o.Settings.NumGoroutines = n
	}
}
