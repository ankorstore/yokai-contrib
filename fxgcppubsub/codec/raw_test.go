package codec_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/stretchr/testify/assert"
)

func TestRawCodec(t *testing.T) {
	t.Parallel()

	t.Run("raw encoding success", func(t *testing.T) {
		t.Parallel()

		rawCodec := codec.NewRawCodec()

		enc, err := rawCodec.Encode([]byte("test"))
		assert.NoError(t, err)
		assert.Equal(t, []byte("test"), enc)

		enc, err = rawCodec.Encode("test")
		assert.NoError(t, err)
		assert.Equal(t, []byte("test"), enc)

		in := struct {
			Test string
		}{
			Test: "test",
		}

		enc, err = rawCodec.Encode(in)
		assert.NoError(t, err)
		assert.Equal(t, []byte(fmt.Sprintf("%s", in)), enc)
	})

	t.Run("raw decoding failure", func(t *testing.T) {
		t.Parallel()

		rawCodec := codec.NewRawCodec()

		err := rawCodec.Decode([]byte("test"), struct{}{})
		assert.Error(t, err)
		assert.Equal(t, "data without schema cannot be decoded", err.Error())
	})
}
