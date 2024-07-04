package subscription

import (
	"time"

	"cloud.google.com/go/pubsub"
)

// Options represents subscription options.
type Options struct {
	ReceiveSettings pubsub.ReceiveSettings
}

// DefaultSubscribeOptions is the default subscription options.
func DefaultSubscribeOptions() *Options {
	return &Options{
		ReceiveSettings: pubsub.DefaultReceiveSettings,
	}
}

// SubscribeOption represents subscription functional options.
type SubscribeOption func(o *Options)

// WithMaxExtension sets the max extension.
func WithMaxExtension(t time.Duration) SubscribeOption {
	return func(o *Options) {
		o.ReceiveSettings.MaxExtension = t
	}
}

// WithMinExtensionPeriod sets the min extension period.
func WithMinExtensionPeriod(t time.Duration) SubscribeOption {
	return func(o *Options) {
		o.ReceiveSettings.MinExtensionPeriod = t
	}
}

// WithMaxExtensionPeriod sets the max extension period.
func WithMaxExtensionPeriod(t time.Duration) SubscribeOption {
	return func(o *Options) {
		o.ReceiveSettings.MaxExtensionPeriod = t
	}
}

// WithMaxOutstandingMessages sets the max outstanding messages.
func WithMaxOutstandingMessages(n int) SubscribeOption {
	return func(o *Options) {
		o.ReceiveSettings.MaxOutstandingMessages = n
	}
}

// WithMaxOutstandingBytes sets the max outstanding bytes.
func WithMaxOutstandingBytes(n int) SubscribeOption {
	return func(o *Options) {
		o.ReceiveSettings.MaxOutstandingBytes = n
	}
}

// WithLegacyFlowControl sets the legacy flow control usage.
func WithLegacyFlowControl(c bool) SubscribeOption {
	return func(o *Options) {
		o.ReceiveSettings.UseLegacyFlowControl = c
	}
}

// WithNumGoroutines sets the num of goroutines.
func WithNumGoroutines(n int) SubscribeOption {
	return func(o *Options) {
		o.ReceiveSettings.NumGoroutines = n
	}
}
