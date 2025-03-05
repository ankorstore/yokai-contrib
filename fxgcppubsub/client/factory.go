package client

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"google.golang.org/api/option"
)

var _ ClientFactory = (*DefaultClientFactory)(nil)

type ClientFactory interface {
	Create(ctx context.Context, projectID string, opts ...option.ClientOption) (*pubsub.Client, error)
}

type DefaultClientFactory struct {
	config *config.Config
	logger *log.Logger
}

func NewDefaultClientFactory(config *config.Config, logger *log.Logger) *DefaultClientFactory {
	return &DefaultClientFactory{
		config: config,
		logger: logger,
	}
}

func (f *DefaultClientFactory) Create(ctx context.Context, projectID string, opts ...option.ClientOption) (*pubsub.Client, error) {
	attempts := f.config.GetInt("modules.gcppubsub.factory.attempts")
	interval := time.Duration(f.config.GetInt("modules.gcppubsub.factory.interval")) * time.Second

	for attempt := 0; attempt <= attempts; attempt++ {
		client, err := pubsub.NewClient(ctx, projectID, opts...)
		if err == nil {
			f.logger.
				Debug().
				Int("attempt", attempt+1).
				Msg("pubsub client creation success")

			return client, nil
		}

		if attempt < attempts {
			f.logger.
				Warn().
				Err(err).
				Int("attempt", attempt+1).
				Msgf("pubsub client creation error, attempting again in %d seconds", int(interval.Seconds()))

			time.Sleep(interval)
		}
	}

	return nil, fmt.Errorf("pubsub client creation error after %d attempts", attempts)
}
