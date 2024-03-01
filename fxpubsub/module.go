package fxpubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"github.com/ankorstore/yokai/config"
	"go.uber.org/fx"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ModuleName is the module name.
const ModuleName = "pubsub"

// FxPubSubModule is the [Fx] pubsub module.
//
// [Fx]: https://github.com/uber-go/fx
var FxPubSubModule = fx.Module(
	ModuleName,
	fx.Provide(
		NewFxPubSub,
	),
)

// FxPubSubParam allows injection of the required dependencies in [NewFxPubSub].
type FxPubSubParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *config.Config
}

// NewFxPubSub returns a [pubsub.Client].
func NewFxPubSub(p FxPubSubParam) (*pubsub.Client, error) {
	var client *pubsub.Client
	var err error

	// client
	if p.Config.IsTestEnv() {
		client, err = createTestClient(p)
	} else {
		client, err = createClient(p)
	}

	return client, err
}

func createClient(p FxPubSubParam) (*pubsub.Client, error) {
	client, err := pubsub.NewClient(context.Background(), p.Config.GetString("modules.pubsub.project.id"))
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %w", err)
	}

	p.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}

func createTestClient(p FxPubSubParam) (*pubsub.Client, error) {
	srv := pstest.NewServer()
	conn, err := grpc.Dial(srv.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client, err := pubsub.NewClient(
		context.Background(),
		p.Config.GetString("modules.pubsub.project.id"),
		option.WithGRPCConn(conn),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create test pubsub client: %w", err)
	}

	p.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err = client.Close()
			if err != nil {
				return err
			}

			return srv.Close()
		},
	})

	return client, nil
}
