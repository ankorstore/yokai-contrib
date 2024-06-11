package fxgcppubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/schema"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/subscription"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/topic"
	"github.com/ankorstore/yokai/config"
	"go.uber.org/fx"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ModuleName is the module name.
const ModuleName = "gcppubsub"

// FxGcpPubSubModule is the [Fx] pubsub module.
//
// [Fx]: https://github.com/uber-go/fx
var FxGcpPubSubModule = fx.Module(
	ModuleName,
	fx.Provide(
		schema.NewSchemaRegistry,
		topic.NewTopicFactory,
		topic.NewTopicRegistry,
		subscription.NewSubscriptionFactory,
		subscription.NewSubscriptionRegistry,
		NewFxGcpPubSubTestServer,
		NewFxGcpPubSubClient,
		NewFxGcpPubSubSchemaClient,
		NewFxGcpPubSubPublisher,
		NewFxGcpPubSubSubscriber,
	),
)

// NewFxGcpPubSubTestServer returns a [pstest.Server].
func NewFxGcpPubSubTestServer() *pstest.Server {
	return pstest.NewServer()
}

// FxGcpPubSubClientParam allows injection of the required dependencies in [NewFxGcpPubSubClient].
//
//nolint:containedctx
type FxGcpPubSubClientParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Context   context.Context
	Config    *config.Config
	Server    *pstest.Server
}

// NewFxGcpPubSubClient returns a [pubsub.Client].
func NewFxGcpPubSubClient(p FxGcpPubSubClientParam) (*pubsub.Client, error) {
	if p.Config.IsTestEnv() {
		client, err := pubsub.NewClient(
			context.Background(),
			p.Config.GetString("modules.gcppubsub.project.id"),
			option.WithEndpoint(p.Server.Addr),
			option.WithoutAuthentication(),
			option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create test pubsub client: %w", err)
		}

		return client, nil
	}

	client, err := pubsub.NewClient(p.Context, p.Config.GetString("modules.gcppubsub.project.id"))
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %w", err)
	}

	p.LifeCycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}

// FxGcpPubSubSchemaClientParam allows injection of the required dependencies in [NewFxGcpPubSubSchemaClient].
//
//nolint:containedctx
type FxGcpPubSubSchemaClientParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Context   context.Context
	Config    *config.Config
	Server    *pstest.Server
}

// NewFxGcpPubSubSchemaClient returns a [pubsub.SchemaClient].
func NewFxGcpPubSubSchemaClient(p FxGcpPubSubSchemaClientParam) (*pubsub.SchemaClient, error) {
	if p.Config.IsTestEnv() {
		client, err := pubsub.NewSchemaClient(
			context.Background(),
			p.Config.GetString("modules.gcppubsub.project.id"),
			option.WithEndpoint(p.Server.Addr),
			option.WithoutAuthentication(),
			option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create test pubsub schema client: %w", err)
		}

		return client, nil
	}

	var schemaClientOptions []option.ClientOption
	if emulatorHost := p.Config.GetEnvVar("PUBSUB_EMULATOR_HOST"); emulatorHost != "" {
		schemaClientOptions = []option.ClientOption{
			option.WithEndpoint(emulatorHost),
			option.WithoutAuthentication(),
			option.WithTelemetryDisabled(),
			internaloption.SkipDialSettingsValidation(),
			option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		}
	}

	client, err := pubsub.NewSchemaClient(context.Background(), p.Config.GetString("modules.gcppubsub.project.id"), schemaClientOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub schema client: %w", err)
	}

	p.LifeCycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}

// FxGcpPubSubPublisherParam allows injection of the required dependencies in [NewFxGcpPubSubPublisher].
type FxGcpPubSubPublisherParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *config.Config
	Factory   *topic.TopicFactory
	Registry  *topic.TopicRegistry
}

// NewFxGcpPubSubPublisher returns a [Publisher].
func NewFxGcpPubSubPublisher(p FxGcpPubSubPublisherParam) *Publisher {
	publisher := NewPublisher(p.Factory, p.Registry)

	if !p.Config.IsTestEnv() {
		p.LifeCycle.Append(fx.Hook{
			OnStop: func(context.Context) error {
				publisher.Stop()

				return nil
			},
		})
	}

	return publisher
}

// FxGcpPubSubSubscriberParam allows injection of the required dependencies in [NewFxGcpPubSubPublisher].
type FxGcpPubSubSubscriberParam struct {
	fx.In
	Factory  *subscription.SubscriptionFactory
	Registry *subscription.SubscriptionRegistry
}

// NewFxGcpPubSubSubscriber returns a [Subscriber].
func NewFxGcpPubSubSubscriber(p FxGcpPubSubSubscriberParam) *Subscriber {
	return NewSubscriber(p.Factory, p.Registry)
}
