package transport

import (
	"net"
	"os"
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

// FindUnusedTestUnixSocketPath returns an unused unix socket path, for tests purposes.
func FindUnusedTestUnixSocketPath(t *testing.T) string {
	t.Helper()

	tempFile, err := os.CreateTemp("", "mysql-*.sock")
	if err != nil {
		t.Errorf("cannot create temp file: %v", err)
	}

	err = tempFile.Close()
	if err != nil {
		t.Errorf("cannot close temp file: %v", err)
	}

	err = os.Remove(tempFile.Name())
	if err != nil {
		t.Errorf("cannot remove temp file: %v", err)
	}

	return tempFile.Name()
}
