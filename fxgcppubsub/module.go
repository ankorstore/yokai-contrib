package fxgcppubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/schema"
	"github.com/ankorstore/yokai/config"
	"go.uber.org/fx"
	"google.golang.org/api/option"
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
		NewFxGcpPubSubClient,
		NewFxGcpPubSubSchemaClient,
		NewFxGcpPubSubSchemaRegistry,
	),
)

// FxGcpPubSubClientParam allows injection of the required dependencies in [NewFxGcpPubSubClient].
type FxGcpPubSubClientParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Context   context.Context
	Config    *config.Config
}

// NewFxGcpPubSubClient returns a [pubsub.Client].
func NewFxGcpPubSubClient(p FxGcpPubSubClientParam) (*pubsub.Client, error) {
	if p.Config.IsTestEnv() {
		return createTestClient(p)
	} else {
		return createClient(p)
	}
}

func createClient(p FxGcpPubSubClientParam) (*pubsub.Client, error) {
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

func createTestClient(p FxGcpPubSubClientParam) (*pubsub.Client, error) {
	srv := pstest.NewServer()

	conn, err := grpc.Dial(srv.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client, err := pubsub.NewClient(
		p.Context,
		p.Config.GetString("modules.gcppubsub.project.id"),
		option.WithGRPCConn(conn),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create test pubsub client: %w", err)
	}

	p.LifeCycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			err = client.Close()
			if err != nil {
				return err
			}

			return srv.Close()
		},
	})

	return client, nil
}

// FxGcpPubSubSchemaClientParam allows injection of the required dependencies in [NewFxGcpPubSubSchemaClient].
type FxGcpPubSubSchemaClientParam struct {
	fx.In
	Context context.Context
	Config  *config.Config
}

// NewFxGcpPubSubSchemaClient returns a [pubsub.SchemaClient].
func NewFxGcpPubSubSchemaClient(p FxGcpPubSubSchemaClientParam) (*pubsub.SchemaClient, error) {
	return pubsub.NewSchemaClient(p.Context, p.Config.GetString("modules.gcppubsub.project.id"))

}

// FxGcpPubSubSchemaRegistryParam allows injection of the required dependencies in [NewFxGcpPubSubSchemaRegistry].
type FxGcpPubSubSchemaRegistryParam struct {
	fx.In
	Client *pubsub.SchemaClient
}

// NewFxGcpPubSubSchemaRegistry returns a [schema.SchemaRegistry].
func NewFxGcpPubSubSchemaRegistry(p FxGcpPubSubSchemaRegistryParam) *schema.SchemaRegistry {
	return schema.NewSchemaRegistry(p.Client)

}
