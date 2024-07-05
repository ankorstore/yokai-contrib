package fxgcppubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor/ack"
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
		fx.Annotate(
			codec.NewDefaultCodecFactory,
			fx.As(new(codec.CodecFactory)),
		),
		fx.Annotate(
			schema.NewDefaultSchemaConfigRegistry,
			fx.As(new(schema.SchemaConfigRegistry)),
		),
		fx.Annotate(
			topic.NewDefaultTopicFactory,
			fx.As(new(topic.TopicFactory)),
		),
		fx.Annotate(
			topic.NewDefaultTopicRegistry,
			fx.As(new(topic.TopicRegistry)),
		),
		fx.Annotate(
			subscription.NewDefaultSubscriptionFactory,
			fx.As(new(subscription.SubscriptionFactory)),
		),
		fx.Annotate(
			subscription.NewDefaultSubscriptionRegistry,
			fx.As(new(subscription.SubscriptionRegistry)),
		),
		fx.Annotate(
			reactor.NewDefaultWaiterSupervisor,
			fx.As(new(reactor.WaiterSupervisor)),
		),
		fx.Annotate(
			NewFxGcpPubSubPublisher,
			fx.As(new(Publisher)),
		),
		fx.Annotate(
			NewFxGcpPubSubSubscriber,
			fx.As(new(Subscriber)),
		),
		NewFxGcpPubSubTestServer,
		NewFxGcpPubSubClient,
		NewFxGcpPubSubSchemaClient,
	),

	AsPubSubTestServerReactor(ack.NewAckReactor),
)

type FxGcpPubSubTestServerParam struct {
	fx.In
	Reactors []Reactor `group:"gcppubsub-reactors"`
}

// NewFxGcpPubSubTestServer returns a [pstest.Server].
func NewFxGcpPubSubTestServer(p FxGcpPubSubTestServerParam) *pstest.Server {
	options := []pstest.ServerReactorOption{}

	for _, r := range p.Reactors {
		for _, fn := range r.FuncNames() {
			options = append(options, pstest.ServerReactorOption{
				FuncName: fn,
				Reactor:  r,
			})
		}
	}

	return pstest.NewServer(options...)
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
	projectID := p.Config.GetString("modules.gcppubsub.project.id")

	if p.Config.IsTestEnv() {
		client, err := pubsub.NewClient(
			p.Context,
			projectID,
			option.WithEndpoint(p.Server.Addr),
			option.WithoutAuthentication(),
			option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create test pubsub client: %w", err)
		}

		return client, nil
	}

	client, err := pubsub.NewClient(p.Context, projectID)
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
	projectID := p.Config.GetString("modules.gcppubsub.project.id")

	if p.Config.IsTestEnv() {
		client, err := pubsub.NewSchemaClient(
			p.Context,
			projectID,
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

	client, err := pubsub.NewSchemaClient(p.Context, projectID, schemaClientOptions...)
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
	Factory   topic.TopicFactory
	Registry  topic.TopicRegistry
}

// NewFxGcpPubSubPublisher returns a [Publisher].
func NewFxGcpPubSubPublisher(p FxGcpPubSubPublisherParam) *DefaultPublisher {
	publisher := NewDefaultPublisher(p.Factory, p.Registry)

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
	Factory  subscription.SubscriptionFactory
	Registry subscription.SubscriptionRegistry
}

// NewFxGcpPubSubSubscriber returns a [Subscriber].
func NewFxGcpPubSubSubscriber(p FxGcpPubSubSubscriberParam) *DefaultSubscriber {
	return NewDefaultSubscriber(p.Factory, p.Registry)
}
