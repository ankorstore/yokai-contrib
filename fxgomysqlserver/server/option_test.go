package server_test

import (
	"io"
	"os"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgomysqlserver/config"
	"github.com/ankorstore/yokai-contrib/fxgomysqlserver/server"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestGoMySQLServerOption(t *testing.T) {
	t.Parallel()

	t.Run("test with default options", func(t *testing.T) {
		t.Parallel()

		options := server.DefaultGoMySQLServerOptions()

		assert.Equal(t, config.TCPTransport, options.Config.Transport())
		assert.Equal(t, config.DefaultUser, options.Config.User())
		assert.Equal(t, config.DefaultPassword, options.Config.Password())
		assert.Equal(t, config.DefaultHost, options.Config.Host())
		assert.Equal(t, config.DefaultPort, options.Config.Port())
		assert.Equal(t, config.DefaultDatabase, options.Config.Database())

		assert.Equal(t, io.Discard, options.LogOutput)

		assert.Equal(t, noop.NewTracerProvider().Tracer("noop"), options.Tracer)
	})

	t.Run("test with custom options", func(t *testing.T) {
		t.Parallel()

		testConfig := config.NewGoMySQLServerConfig(config.WithTransport(config.MemoryTransport))
		testLogOutput := os.Stdout
		testTracer := otel.GetTracerProvider().Tracer("test")

		options := &server.ServerOptions{}

		server.WithConfig(testConfig)(options)
		server.WithLogOutput(testLogOutput)(options)
		server.WithTracer(testTracer)(options)

		assert.Equal(t, testConfig, options.Config)
		assert.Equal(t, testLogOutput, options.LogOutput)
		assert.Equal(t, testTracer, options.Tracer)
	})
}
