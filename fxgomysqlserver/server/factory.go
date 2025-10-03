package server

import (
	"context"
	"fmt"
	"net"

	"github.com/ankorstore/yokai-contrib/fxgomysqlserver/config"
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var _ GoMySQLServerFactory = (*DefaultGoMySQLServerFactory)(nil)

type GoMySQLServerFactory interface {
	Create(options ...GoMySQLServerOption) (*server.Server, error)
}

type DefaultGoMySQLServerFactory struct{}

func NewDefaultGoMySQLServerFactory() *DefaultGoMySQLServerFactory {
	return &DefaultGoMySQLServerFactory{}
}

func (f *DefaultGoMySQLServerFactory) Create(options ...GoMySQLServerOption) (*server.Server, error) {
	// server options
	serverOptions := DefaultGoMySQLServerOptions()
	for _, opt := range options {
		opt(&serverOptions)
	}

	// server logger output
	logrus.SetOutput(serverOptions.LogOutput)

	// server engine
	serverDB := memory.NewDatabase(serverOptions.Config.Database())
	serverDB.BaseDatabase.EnablePrimaryKeyIndexes()
	serverProvider := memory.NewDBProvider(serverDB)
	serverEngine := sqle.NewDefault(serverProvider)

	// server superuser
	serverCatalog := serverEngine.Analyzer.Catalog.MySQLDb
	serverCatalogEditor := serverCatalog.Editor()
	defer serverCatalogEditor.Close()

	serverCatalog.AddSuperUser(
		serverCatalogEditor,
		serverOptions.Config.User(),
		serverOptions.Config.Host(),
		serverOptions.Config.Password(),
	)

	// server session builder
	serverSessionBuilder := memory.NewSessionBuilder(serverProvider)

	// server config
	serverConfig := server.Config{
		Tracer: serverOptions.Tracer,
	}

	// transport specific server config
	switch serverTransport := serverOptions.Config.Transport(); serverTransport {
	case config.TCPTransport:
		serverConfig.Protocol = "tcp"
		serverConfig.Address = fmt.Sprintf("%s:%d", serverOptions.Config.Host(), serverOptions.Config.Port())
	case config.MemoryTransport:
		serverConfig.Protocol = "memory"
		serverConfig.Listener = serverOptions.Config.Listener()

		// register listener in mysql driver
		mysql.RegisterDialContext(config.DefaultMemoryNetwork, func(ctx context.Context, addr string) (net.Conn, error) {
			return serverOptions.Config.Listener().DialContext(ctx)
		})

		// activate superuser also on memory listener network
		serverCatalog.AddSuperUser(
			serverCatalogEditor,
			serverOptions.Config.User(),
			config.DefaultMemoryAddress,
			serverOptions.Config.Password(),
		)
	case config.UnknownTransport:
		return nil, fmt.Errorf("unknown transport: %s", serverTransport.String())
	}

	// create server
	return server.NewServer(serverConfig, serverEngine, sql.NewContext, serverSessionBuilder, nil)
}
