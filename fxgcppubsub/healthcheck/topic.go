package healthcheck

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/healthcheck"
)

// TopicsProbeName is the name of the GCP pub/sub topics probe.
const TopicsProbeName = "gcppubsub-topics"

// GcpPubSubTopicsProbe is a probe compatible with the [healthcheck] module.
//
// [healthcheck]: https://github.com/ankorstore/yokai/tree/main/healthcheck
type GcpPubSubTopicsProbe struct {
	config *config.Config
	client *pubsub.Client
}

// NewGcpPubSubTopicsProbe returns a new [GcpPubSubTopicsProbe].
func NewGcpPubSubTopicsProbe(config *config.Config, client *pubsub.Client) *GcpPubSubTopicsProbe {
	return &GcpPubSubTopicsProbe{
		config: config,
		client: client,
	}
}

// Name returns the name of the [GcpPubSubTopicsProbe].
func (p *GcpPubSubTopicsProbe) Name() string {
	return TopicsProbeName
}

// Check returns a successful [healthcheck.CheckerProbeResult] if all configured topics exist.
func (p *GcpPubSubTopicsProbe) Check(ctx context.Context) *healthcheck.CheckerProbeResult {
	success := true
	var messages []string

	for _, topicName := range p.config.GetStringSlice("modules.gcppubsub.healthcheck.topics") {
		topic := p.client.Topic(topicName)

		exists, err := topic.Exists(ctx)
		if err != nil {
			success = false
			messages = append(messages, fmt.Sprintf("topic %s error: %v", topicName, err))
		} else {
			if !exists {
				success = false
				messages = append(messages, fmt.Sprintf("topic %s does not exist", topicName))
			} else {
				messages = append(messages, fmt.Sprintf("topic %s exists", topicName))
			}
		}
	}

	return healthcheck.NewCheckerProbeResult(success, strings.Join(messages, ", "))
}
