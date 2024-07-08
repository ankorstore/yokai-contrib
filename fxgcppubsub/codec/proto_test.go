package codec_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/proto"
	"github.com/stretchr/testify/assert"
)

func TestProtoBinaryCodec(t *testing.T) {
	t.Parallel()

	schemaDefinition := proto.GetTestProtoSchemaDefinition(t)

	t.Run("protobuf binary encoding and decoding success", func(t *testing.T) {
		t.Parallel()

		protoBinaryCodec := codec.NewProtoBinaryCodec(schemaDefinition)

		in := &proto.SimpleRecord{
			StringField:  "test",
			FloatField:   12.34,
			BooleanField: true,
		}

		enc, err := protoBinaryCodec.Encode(in)
		assert.NoError(t, err)

		out := proto.SimpleRecord{}

		err = protoBinaryCodec.Decode(enc, &out)
		assert.NoError(t, err)

		assert.Equal(t, in.StringField, out.StringField)
		assert.Equal(t, in.FloatField, out.FloatField)
		assert.Equal(t, in.BooleanField, out.BooleanField)
	})

	t.Run("protobuf binary encoding failure", func(t *testing.T) {
		t.Parallel()

		protoBinaryCodec := codec.NewProtoBinaryCodec(schemaDefinition)

		_, err := protoBinaryCodec.Encode(struct{}{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid proto message")
	})

	t.Run("protobuf binary decoding failure", func(t *testing.T) {
		t.Parallel()

		protoBinaryCodec := codec.NewProtoBinaryCodec(schemaDefinition)

		out := proto.SimpleRecord{}

		err := protoBinaryCodec.Decode([]byte("invalid"), &out)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot decode proto binary")
	})
}

func TestProtoJsonCodec(t *testing.T) {
	t.Parallel()

	schemaDefinition := proto.GetTestProtoSchemaDefinition(t)

	t.Run("protobuf json encoding and decoding success", func(t *testing.T) {
		t.Parallel()

		protoJsonCodec := codec.NewProtoJsonCodec(schemaDefinition)

		in := &proto.SimpleRecord{
			StringField:  "test",
			FloatField:   12.34,
			BooleanField: true,
		}

		enc, err := protoJsonCodec.Encode(in)
		assert.NoError(t, err)

		out := proto.SimpleRecord{}

		err = protoJsonCodec.Decode(enc, &out)
		assert.NoError(t, err)

		assert.Equal(t, in.StringField, out.StringField)
		assert.Equal(t, in.FloatField, out.FloatField)
		assert.Equal(t, in.BooleanField, out.BooleanField)
	})

	t.Run("protobuf json encoding failure", func(t *testing.T) {
		t.Parallel()

		protoJsonCodec := codec.NewProtoJsonCodec(schemaDefinition)

		_, err := protoJsonCodec.Encode(struct{}{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid proto message")
	})

	t.Run("protobuf json decoding failure", func(t *testing.T) {
		t.Parallel()

		protoJsonCodec := codec.NewProtoJsonCodec(schemaDefinition)

		out := proto.SimpleRecord{}

		err := protoJsonCodec.Decode([]byte("invalid"), &out)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot decode proto json")
	})
}
