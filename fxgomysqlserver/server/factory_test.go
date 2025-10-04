package server_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgomysqlserver/config"
	"github.com/ankorstore/yokai-contrib/fxgomysqlserver/server"
	"github.com/ankorstore/yokai-contrib/fxgomysqlserver/testdata/transport"
	sqle "github.com/dolthub/go-mysql-server/server"
	"github.com/stretchr/testify/assert"
)

func TestDefaultGoMySQLServerFactory(t *testing.T) {
	t.Run("test implementation", func(t *testing.T) {
		serverFactory := server.NewDefaultGoMySQLServerFactory()

		assert.IsType(t, &server.DefaultGoMySQLServerFactory{}, serverFactory)
		assert.Implements(t, (*server.GoMySQLServerFactory)(nil), serverFactory)
	})

	t.Run("test tcp server creation", func(t *testing.T) {
		serverPort := transport.FindUnusedTestTCPPort(t)

		serverConfig := config.NewGoMySQLServerConfig(
			config.WithPort(serverPort),
		)

		srv, err := server.NewDefaultGoMySQLServerFactory().Create(
			server.WithConfig(serverConfig),
		)
		assert.NoError(t, err)

		assert.IsType(t, &sqle.Server{}, srv)
		assert.Equal(t, "tcp", srv.Listener.Addr().Network())
		assert.Equal(t, fmt.Sprintf("127.0.0.1:%d", serverPort), srv.Listener.Addr().String())

		srv.Listener.Close()
	})

	t.Run("test memory server creation", func(t *testing.T) {
		serverConfig := config.NewGoMySQLServerConfig(
			config.WithTransport(config.MemoryTransport),
		)

		srv, err := server.NewDefaultGoMySQLServerFactory().Create(
			server.WithConfig(serverConfig),
		)
		assert.NoError(t, err)

		assert.IsType(t, &sqle.Server{}, srv)
		assert.Equal(t, config.DefaultMemoryAddress, srv.Listener.Addr().String())

		srv.Listener.Close()
	})

	t.Run("test creation failure with invalid transport", func(t *testing.T) {
		serverConfig := config.NewGoMySQLServerConfig(
			config.WithTransport(config.UnknownTransport),
		)

		srv, err := server.NewDefaultGoMySQLServerFactory().Create(
			server.WithConfig(serverConfig),
		)
		assert.Nil(t, srv)
		assert.Error(t, err)
		assert.Equal(t, "unknown transport: unknown", err.Error())
	})
}
