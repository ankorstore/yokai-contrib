package proto

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed simple.proto
var contents []byte

func GetTestProtoSchemaDefinition(tb testing.TB) string {
	tb.Helper()

	assert.NotEmpty(tb, contents)

	return string(contents)
}
