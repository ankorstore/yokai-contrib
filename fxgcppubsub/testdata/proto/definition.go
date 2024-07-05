package proto

import (
	_ "embed"
	"testing"
)

//go:embed simple.proto
var contents []byte

func GetTestProtoSchemaDefinition(tb testing.TB) string {
	tb.Helper()

	return string(contents)
}
