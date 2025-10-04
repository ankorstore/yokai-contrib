package transport_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgomysqlserver/testdata/transport"
	"github.com/stretchr/testify/assert"
)

func TestFindUnusedTestTCPPort(t *testing.T) {
	t.Parallel()

	port1 := transport.FindUnusedTestTCPPort(t)
	port2 := transport.FindUnusedTestTCPPort(t)
	port3 := transport.FindUnusedTestTCPPort(t)

	assert.NotEqual(t, port1, port2)
	assert.NotEqual(t, port2, port3)
	assert.NotEqual(t, port1, port3)
}
