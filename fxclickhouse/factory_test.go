package fxclickhouse_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxclickhouse"
	"github.com/stretchr/testify/assert"
)

func TestDefaultClickhouseFactory(t *testing.T) {
	t.Parallel()

	factory := fxclickhouse.NewDefaultClickhouseFactory()

	assert.IsType(t, &fxclickhouse.DefaultClickhouseFactory{}, factory)
	assert.Implements(t, (*fxclickhouse.ClickhouseFactory)(nil), factory)
}
