package codec_test

import (
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/stretchr/testify/assert"
)

func TestDefaultCodecFactory(t *testing.T) {
	t.Parallel()

	t.Run("construction", func(t *testing.T) {
		t.Parallel()

		defaultFactory := codec.NewDefaultCodecFactory()

		assert.IsType(t, &codec.DefaultCodecFactory{}, defaultFactory)
		assert.Implements(t, (*codec.CodecFactory)(nil), defaultFactory)
	})

	t.Run("codec creation", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodecFactory().Create(
			pubsub.SchemaTypeUnspecified,
			pubsub.EncodingUnspecified,
			"",
		)

		assert.IsType(t, &codec.DefaultCodec{}, defaultCodec)
		assert.Implements(t, (*codec.Codec)(nil), defaultCodec)
	})
}
