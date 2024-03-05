package healthcheck

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/healthcheck"
)

// SubscriptionsProbeName is the name of the GCP pub/sub subscriptions probe.
const SubscriptionsProbeName = "gcppubsub-subscriptions"

// GcpPubSubSubscriptionsProbe is a probe compatible with the [healthcheck] module.
//
// [healthcheck]: https://github.com/ankorstore/yokai/tree/main/healthcheck
type GcpPubSubSubscriptionsProbe struct {
	config *config.Config
	client *pubsub.Client
}

// NewGcpPubSubSubscriptionsProbe returns a new [GcpPubSubSubscriptionsProbe].
func NewGcpPubSubSubscriptionsProbe(config *config.Config, client *pubsub.Client) *GcpPubSubSubscriptionsProbe {
	return &GcpPubSubSubscriptionsProbe{
		config: config,
		client: client,
	}
}

// Name returns the name of the [GcpPubSubSubscriptionsProbe].
func (p *GcpPubSubSubscriptionsProbe) Name() string {
	return SubscriptionsProbeName
}

// Check returns a successful [healthcheck.CheckerProbeResult] if all configured subscriptions exist.
func (p *GcpPubSubSubscriptionsProbe) Check(ctx context.Context) *healthcheck.CheckerProbeResult {
	success := true
	var messages []string

	for _, subscriptionName := range p.config.GetStringSlice("modules.gcppubsub.healthcheck.subscriptions") {
		subscription := p.client.Subscription(subscriptionName)

		exists, err := subscription.Exists(ctx)
		if err != nil {
			success = false
			messages = append(messages, fmt.Sprintf("subscription %s error: %v", subscriptionName, err))
		} else {
			if !exists {
				success = false
				messages = append(messages, fmt.Sprintf("subscription %s does not exist", subscriptionName))
			} else {
				messages = append(messages, fmt.Sprintf("subscription %s exists", subscriptionName))
			}
		}
	}

	return healthcheck.NewCheckerProbeResult(success, strings.Join(messages, ", "))
}
