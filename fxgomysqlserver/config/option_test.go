package config_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgomysqlserver/config"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/test/bufconn"
)

func TestGoMySQLServerConfigOption(t *testing.T) {
	t.Parallel()

	t.Run("test with default options", func(t *testing.T) {
		t.Parallel()

		options := config.DefaultGoMySQLServerConfigOptions()

		assert.Equal(t, config.TCPTransport, options.Transport)
		assert.Equal(t, config.DefaultUser, options.User)
		assert.Equal(t, config.DefaultPassword, options.Password)
		assert.Equal(t, config.DefaultSocket, options.Socket)
		assert.Equal(t, config.DefaultHost, options.Host)
		assert.Equal(t, config.DefaultPort, options.Port)
		assert.Equal(t, config.DefaultDatabase, options.Database)

		assert.IsType(t, &bufconn.Listener{}, options.Listener)
	})

	t.Run("test with custom options", func(t *testing.T) {
		t.Parallel()

		options := &config.ConfigOptions{}

		config.WithTransport(config.SocketTransport)(options)
		config.WithUser("test_user")(options)
		config.WithPassword("test_password")(options)
		config.WithSocket("test_socket")(options)
		config.WithHost("test_host")(options)
		config.WithPort(9999)(options)
		config.WithDatabase("test_database")(options)

		assert.Equal(t, config.SocketTransport, options.Transport)
		assert.Equal(t, "test_user", options.User)
		assert.Equal(t, "test_password", options.Password)
		assert.Equal(t, "test_socket", options.Socket)
		assert.Equal(t, "test_host", options.Host)
		assert.Equal(t, 9999, options.Port)
		assert.Equal(t, "test_database", options.Database)

		assert.IsType(t, &bufconn.Listener{}, options.Listener)
	})
}
