package message_test

import (
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/message"
	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	t.Parallel()

	t.Run("base message wrapping", func(t *testing.T) {
		t.Parallel()

		baseMsg := createTestBaseMessage()
		cod := codec.NewRawCodec()

		msg := message.NewMessage(cod, baseMsg)

		assert.Equal(t, baseMsg, msg.BaseMessage())
		assert.Equal(t, cod, msg.Codec())

		assert.Equal(t, "foo", msg.ID())
		assert.Equal(t, []byte("bar"), msg.Data())
		assert.Equal(t, map[string]string{"baz": "baz"}, msg.Attributes())
	})

	t.Run("message decoding failure without schema", func(t *testing.T) {
		t.Parallel()

		baseMsg := createTestBaseMessage()
		cod := codec.NewRawCodec()

		msg := message.NewMessage(cod, baseMsg)

		var out []byte
		err := msg.Decode(&out)
		assert.Error(t, err)
		assert.Equal(t, "data without schema cannot be decoded", err.Error())
	})
}

func createTestBaseMessage() *pubsub.Message {
	return &pubsub.Message{
		ID:         "foo",
		Data:       []byte("bar"),
		Attributes: map[string]string{"baz": "baz"},
	}
}
