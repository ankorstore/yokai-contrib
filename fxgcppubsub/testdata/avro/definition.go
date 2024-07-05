package avro

import (
	_ "embed"
	"testing"
)

//go:embed simple.avsc
var contents []byte

func GetTestAvroSchemaDefinition(tb testing.TB) string {
	tb.Helper()

	return string(contents)
}
