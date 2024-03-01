package healthcheck

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/healthcheck"
)

type PubSubProbe struct {
	config *config.Config
	client *pubsub.Client
}

func NewPubSubProbe(config *config.Config, client *pubsub.Client) *PubSubProbe {
	return &PubSubProbe{
		config: config,
		client: client,
	}
}

func (p *PubSubProbe) Name() string {
	return "pubsub"
}

func (p *PubSubProbe) Check(ctx context.Context) *healthcheck.CheckerProbeResult {
	topicName := p.config.GetString("modules.pubsub.healthcheck.topics")
	topic := p.client.Topic(topicName)

	exists, err := topic.Exists(ctx)

	if err != nil {
		return healthcheck.NewCheckerProbeResult(false, "pubsub unreachable")
	}

	if !exists {
		return healthcheck.NewCheckerProbeResult(false, fmt.Sprintf("error: topic %s does not exist", topicName))
	}

	return healthcheck.NewCheckerProbeResult(true, fmt.Sprintf("success: topic %s exists", topicName))
}
