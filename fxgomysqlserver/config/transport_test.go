package config_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgomysqlserver/config"
	"github.com/stretchr/testify/assert"
)

func TestTransport(t *testing.T) {
	t.Parallel()

	t.Run("test transport as string", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			transport config.Transport
			expected  string
		}{
			{
				config.TCPTransport,
				"tcp",
			},
			{
				config.SocketTransport,
				"socket",
			},
			{
				config.MemoryTransport,
				"memory",
			},
			{
				config.UnknownTransport,
				"unknown",
			},
		}

		for _, test := range tests {
			assert.Equal(t, test.expected, test.transport.String())
		}
	})

	t.Run("test fetch transport", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			transport string
			expected  config.Transport
		}{
			{
				"tcp",
				config.TCPTransport,
			},
			{
				"TCP",
				config.TCPTransport,
			},
			{
				"Tcp",
				config.TCPTransport,
			},
			{
				"socket",
				config.SocketTransport,
			},
			{
				"SOCKET",
				config.SocketTransport,
			},
			{
				"Socket",
				config.SocketTransport,
			},
			{
				"memory",
				config.MemoryTransport,
			},
			{
				"MEMORY",
				config.MemoryTransport,
			},
			{
				"Memory",
				config.MemoryTransport,
			},
			{
				"",
				config.UnknownTransport,
			},
			{
				"unknown",
				config.UnknownTransport,
			},
			{
				"invalid",
				config.UnknownTransport,
			},
		}

		for _, test := range tests {
			assert.Equal(t, test.expected, config.FetchTransport(test.transport))
		}
	})
}
