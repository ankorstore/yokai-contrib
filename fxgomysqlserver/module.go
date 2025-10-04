package fxgomysqlserver

import (
	"context"
	"fmt"
	"os"

	"github.com/ankorstore/yokai-contrib/fxgomysqlserver/config"
	"github.com/ankorstore/yokai-contrib/fxgomysqlserver/server"
	yokaiconfig "github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	sqle "github.com/dolthub/go-mysql-server/server"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "gomysqlserver"

// FxGoMySQLServerModule is the [Fx] go-mysql-server module.
//
// [Fx]: https://github.com/uber-go/fx
var FxGoMySQLServerModule = fx.Module(
	ModuleName,
	fx.Provide(
		fx.Annotate(
			server.NewDefaultGoMySQLServerFactory,
			fx.As(new(server.GoMySQLServerFactory)),
		),
		fx.Annotate(
			NewFxGoMySQLServerModuleInfo,
			fx.As(new(interface{})),
			fx.ResultTags(`group:"core-module-infos"`),
		),
		NewFxGoMySQLServerConfig,
		NewFxGoMySQLServer,
	),
)

// FxGoMySQLServerConfigParam allows injection of the required dependencies in [NewGoMySQLServerConfig].
type FxGoMySQLServerConfigParam struct {
	fx.In
	Config *yokaiconfig.Config
}

// NewFxGoMySQLServerConfig returns a new [config.GoMySQLServerConfig] instance.
func NewFxGoMySQLServerConfig(p FxGoMySQLServerConfigParam) (*config.GoMySQLServerConfig, error) {
	// server transport
	transportConfig := p.Config.GetString("modules.gomysqlserver.config.transport")
	transport := config.FetchTransport(transportConfig)
	if transport == config.UnknownTransport {
		return nil, fmt.Errorf("unknown transport: %s", transportConfig)
	}

	// server config
	return config.NewGoMySQLServerConfig(
		config.WithTransport(transport),
		config.WithUser(p.Config.GetString("modules.gomysqlserver.config.user")),
		config.WithPassword(p.Config.GetString("modules.gomysqlserver.config.password")),
		config.WithHost(p.Config.GetString("modules.gomysqlserver.config.host")),
		config.WithPort(p.Config.GetInt("modules.gomysqlserver.config.port")),
		config.WithDatabase(p.Config.GetString("modules.gomysqlserver.config.database")),
	), nil
}

// FxGoMySQLServerParam allows injection of the required dependencies in [NewFxGoMySQLServer].
type FxGoMySQLServerParam struct {
	fx.In
	LifeCycle      fx.Lifecycle
	ServerFactory  server.GoMySQLServerFactory
	ServerConfig   *config.GoMySQLServerConfig
	Config         *yokaiconfig.Config
	Logger         *log.Logger
	TracerProvider oteltrace.TracerProvider
}

// NewFxGoMySQLServer returns a new [sqle.Server] instance.
func NewFxGoMySQLServer(p FxGoMySQLServerParam) (*sqle.Server, error) {
	// server config
	options := []server.GoMySQLServerOption{
		server.WithConfig(p.ServerConfig),
	}

	// server logger
	if p.Config.GetBool("modules.gomysqlserver.log.enabled") {
		options = append(options, server.WithLogOutput(os.Stdout))
	}

	// server tracer
	if p.Config.GetBool("modules.gomysqlserver.trace.enabled") {
		options = append(options, server.WithTracer(p.TracerProvider.Tracer(ModuleName)))
	}

	// server creation
	srv, err := p.ServerFactory.Create(options...)
	if err != nil {
		return nil, err
	}

	// server dsn
	dsn, err := p.ServerConfig.DSN()
	if err != nil {
		return nil, err
	}

	// server start
	p.Logger.Info().Msgf("starting go mysql server, accepting connections on %s", dsn)
	//nolint:errcheck
	go srv.Start()

	// server lifecycle
	p.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if p.ServerConfig.Transport() != config.MemoryTransport {
				p.Logger.Info().Msg("stopping go mysql server")

				return srv.Close()
			}

			return nil
		},
	})

	return srv, nil
}
