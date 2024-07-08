package avro

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed simple.avsc
var contents []byte

func GetTestAvroSchemaDefinition(tb testing.TB) string {
	tb.Helper()

	assert.NotEmpty(tb, contents)

	return string(contents)
}
