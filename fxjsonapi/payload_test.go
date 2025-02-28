package fxjsonapi_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxjsonapi"
	"github.com/ankorstore/yokai-contrib/fxjsonapi/testdata/model"
	"github.com/stretchr/testify/assert"
)

//nolint:goconst
func TestMarshall(t *testing.T) {
	t.Parallel()

	t.Run("test failure with invalid data", func(t *testing.T) {
		t.Parallel()

		invalid := struct {
			ID int `jsonapi:"invalid"`
		}{}

		_, err := fxjsonapi.Marshall(&invalid, fxjsonapi.MarshallParams{})
		assert.Error(t, err)
		assert.Equal(t, "Bad jsonapi struct tag format", err.Error())
	})

	t.Run("test success with defaults", func(t *testing.T) {
		t.Parallel()

		foo := model.CreateTestFoo()

		mFoo, err := fxjsonapi.Marshall(&foo, fxjsonapi.MarshallParams{})
		assert.NoError(t, err)

		expected := `{"data":{"type":"foo","id":"123","attributes":{"name":"foo"},"relationships":{"bar":{"data":{"type":"bar","id":"456"}}},"meta":{"meta":"foo"}},"included":[{"type":"bar","id":"456","attributes":{"name":"bar"},"meta":{"meta":"bar"}}]}`

		assert.Equal(t, fmt.Sprintf("%s\n", expected), string(mFoo))
	})

	t.Run("test success with metadata", func(t *testing.T) {
		t.Parallel()

		foo := model.CreateTestFoo()

		mFoo, err := fxjsonapi.Marshall(&foo, fxjsonapi.MarshallParams{
			Metadata: map[string]interface{}{"foo": "bar"},
		})
		assert.NoError(t, err)

		expected := `{"data":{"type":"foo","id":"123","attributes":{"name":"foo"},"relationships":{"bar":{"data":{"type":"bar","id":"456"}}},"meta":{"meta":"foo"}},"included":[{"type":"bar","id":"456","attributes":{"name":"bar"},"meta":{"meta":"bar"}}],"meta":{"foo":"bar"}}`

		assert.Equal(t, fmt.Sprintf("%s\n", expected), string(mFoo))
	})

	t.Run("test success without included", func(t *testing.T) {
		t.Parallel()

		foo := model.CreateTestFoo()

		mFoo, err := fxjsonapi.Marshall(&foo, fxjsonapi.MarshallParams{
			WithoutIncluded: true,
		})
		assert.NoError(t, err)

		expected := `{"data":{"type":"foo","id":"123","attributes":{"name":"foo"},"relationships":{"bar":{"data":{"type":"bar","id":"456"}}},"meta":{"meta":"foo"}}}`

		assert.Equal(t, fmt.Sprintf("%s\n", expected), string(mFoo))
	})
}
