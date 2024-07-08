package subscription

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/message"
)

// SubscribeFunc represents the Subscription execution callback.
type SubscribeFunc func(ctx context.Context, m *message.Message)

// Subscription represents a pub/sub subscription with an associated codec.Codec.
type Subscription struct {
	codec        codec.Codec
	subscription *pubsub.Subscription
	options      *Options
}

// NewSubscription returns a new Subscription instance.
func NewSubscription(codec codec.Codec, subscription *pubsub.Subscription) *Subscription {
	return &Subscription{
		codec:        codec,
		subscription: subscription,
		options:      DefaultSubscribeOptions(),
	}
}

// Codec returns the subscription associated codec.Codec.
func (s *Subscription) Codec() codec.Codec {
	return s.codec
}

// BaseSubscription returns the base pubsub.Subscription.
func (s *Subscription) BaseSubscription() *pubsub.Subscription {
	return s.subscription
}

// WithOptions configures the subscription with a list of SubscribeOption.
func (s *Subscription) WithOptions(options ...SubscribeOption) *Subscription {
	// resolve options
	for _, applyOpt := range options {
		applyOpt(s.options)
	}

	// apply options
	s.subscription.ReceiveSettings = s.options.ReceiveSettings

	return s
}

// Subscribe starts the subscription and runs the provided SubscribeFunc.
func (s *Subscription) Subscribe(ctx context.Context, f SubscribeFunc) error {
	return s.subscription.Receive(ctx, func(fCtx context.Context, msg *pubsub.Message) {
		f(fCtx, message.NewMessage(s.codec, msg))
	})
}
