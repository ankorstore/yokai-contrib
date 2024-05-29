package config_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgomysqlserver/config"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/test/bufconn"
)

func TestGoMySQLServerConfig(t *testing.T) {
	t.Parallel()

	t.Run("test with default tcp options", func(t *testing.T) {
		t.Parallel()

		configuration := config.NewGoMySQLServerConfig()

		assert.Equal(t, config.TCPTransport, configuration.Transport())
		assert.Equal(t, config.DefaultUser, configuration.User())
		assert.Equal(t, config.DefaultPassword, configuration.Password())
		assert.Equal(t, config.DefaultSocket, configuration.Socket())
		assert.Equal(t, config.DefaultHost, configuration.Host())
		assert.Equal(t, config.DefaultPort, configuration.Port())
		assert.Equal(t, config.DefaultDatabase, configuration.Database())

		assert.IsType(t, &bufconn.Listener{}, configuration.Listener())

		dsn, err := configuration.DSN()
		assert.NoError(t, err)
		assert.Equal(t, "user:password@tcp(localhost:3306)/db", dsn)
	})

	t.Run("test with socket options", func(t *testing.T) {
		t.Parallel()

		configuration := config.NewGoMySQLServerConfig(
			config.WithTransport(config.SocketTransport),
			config.WithUser("test_socket_user"),
			config.WithPassword("test_socket_password"),
			config.WithSocket("test_socket_socket"),
			config.WithDatabase("test_socket_database"),
		)

		assert.Equal(t, config.SocketTransport, configuration.Transport())
		assert.Equal(t, "test_socket_user", configuration.User())
		assert.Equal(t, "test_socket_password", configuration.Password())
		assert.Equal(t, "test_socket_socket", configuration.Socket())
		assert.Equal(t, "test_socket_database", configuration.Database())

		assert.IsType(t, &bufconn.Listener{}, configuration.Listener())

		dsn, err := configuration.DSN()
		assert.NoError(t, err)
		assert.Equal(t, "test_socket_user:test_socket_password@unix(test_socket_socket)/test_socket_database", dsn)
	})

	t.Run("test with memory options", func(t *testing.T) {
		t.Parallel()

		configuration := config.NewGoMySQLServerConfig(
			config.WithTransport(config.MemoryTransport),
			config.WithUser("test_memory_user"),
			config.WithPassword("test_memory_password"),
			config.WithDatabase("test_memory_database"),
		)

		assert.Equal(t, config.MemoryTransport, configuration.Transport())
		assert.Equal(t, "test_memory_user", configuration.User())
		assert.Equal(t, "test_memory_password", configuration.Password())
		assert.Equal(t, "test_memory_database", configuration.Database())

		assert.IsType(t, &bufconn.Listener{}, configuration.Listener())

		dsn, err := configuration.DSN()
		assert.NoError(t, err)
		assert.Equal(t, "test_memory_user:test_memory_password@memory(bufconn)/test_memory_database", dsn)
	})

	t.Run("test with invalid transport", func(t *testing.T) {
		t.Parallel()

		configuration := config.NewGoMySQLServerConfig(
			config.WithTransport(config.UnknownTransport),
		)

		_, err := configuration.DSN()
		assert.Error(t, err)
		assert.Equal(t, "unknown transport: unknown", err.Error())
	})
}
