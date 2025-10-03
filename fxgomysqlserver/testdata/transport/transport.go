package transport

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

// FindUnusedTestTCPPort returns an unused TCP port, for tests purposes.
func FindUnusedTestTCPPort(t *testing.T) int {
	t.Helper()

	tempListener, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)

	tempListenerAddr, ok := tempListener.Addr().(*net.TCPAddr)
	if !ok {
		t.Error("temp listener address is not TCP")
	}

	port := tempListenerAddr.Port
	assert.NoError(t, tempListener.Close())

	return port
}
