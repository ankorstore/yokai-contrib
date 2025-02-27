package fxjsonapi_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxjsonapi"
	"github.com/ankorstore/yokai/config"
	"github.com/stretchr/testify/assert"
)

func TestProcessorOptions(t *testing.T) {
	t.Parallel()

	cfg, err := config.NewDefaultConfigFactory().Create(config.WithFilePaths("./testdata/config"))
	assert.NoError(t, err)

	options := fxjsonapi.DefaultProcessorOptions(cfg)

	t.Run("test defaults", func(t *testing.T) {
		t.Parallel()

		assert.Len(t, options.Metadata, 0)
		assert.True(t, options.Included)
		assert.True(t, options.Log)
		assert.True(t, options.Trace)
	})

	t.Run("test with metadata", func(t *testing.T) {
		t.Parallel()

		opt := fxjsonapi.WithMetadata(map[string]any{"foo": "bar"})
		opt(&options)

		assert.Equal(t, "bar", options.Metadata["foo"])
	})

	t.Run("test with included", func(t *testing.T) {
		t.Parallel()

		opt := fxjsonapi.WithIncluded(false)
		opt(&options)

		assert.False(t, options.Included)
	})

	t.Run("test with log", func(t *testing.T) {
		t.Parallel()

		opt := fxjsonapi.WithLog(true)
		opt(&options)

		assert.True(t, options.Log)
	})

	t.Run("test with trace", func(t *testing.T) {
		t.Parallel()

		opt := fxjsonapi.WithTrace(true)
		opt(&options)

		assert.True(t, options.Trace)
	})
}
