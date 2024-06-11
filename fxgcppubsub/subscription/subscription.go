package subscription

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/message"
)

type SubscribeFunc func(ctx context.Context, m *message.Message)

type Subscription struct {
	subscription *pubsub.Subscription
	codec        *codec.Codec
}

func NewSubscription(subscription *pubsub.Subscription, codec *codec.Codec) *Subscription {
	return &Subscription{
		subscription: subscription,
		codec:        codec,
	}
}

func (s *Subscription) Base() *pubsub.Subscription {
	return s.subscription
}

func (s *Subscription) WithOptions(options ...SubscribeOption) *Subscription {
	// resolve options
	subscriptionOptions := DefaultSubscribeOptions()
	for _, applyOpt := range options {
		applyOpt(&subscriptionOptions)
	}

	// apply options
	s.subscription.ReceiveSettings = subscriptionOptions.Settings

	return s
}

func (s *Subscription) Subscribe(ctx context.Context, f SubscribeFunc) error {
	return s.subscription.Receive(ctx, func(c context.Context, msg *pubsub.Message) {
		message := message.NewMessage(msg, s.codec)

		f(ctx, message)
	})
}
