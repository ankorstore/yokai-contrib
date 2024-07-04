package codec_test

import (
	"fmt"
	"os"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/avro"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/proto"
	"github.com/stretchr/testify/assert"
)

//nolint:maintidx
func TestDefaultCodec(t *testing.T) {
	t.Parallel()

	t.Run("raw encoding success", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaTypeUnspecified, pubsub.EncodingUnspecified, "")

		enc, err := defaultCodec.Encode([]byte("test"))
		assert.NoError(t, err)
		assert.Equal(t, []byte("test"), enc)

		enc, err = defaultCodec.Encode("test")
		assert.NoError(t, err)
		assert.Equal(t, []byte("test"), enc)

		in := struct {
			Test string
		}{
			Test: "test",
		}

		enc, err = defaultCodec.Encode(in)
		assert.NoError(t, err)
		assert.Equal(t, []byte(fmt.Sprintf("%s", in)), enc)
	})

	t.Run("raw decoding failure", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaTypeUnspecified, pubsub.EncodingUnspecified, "")

		err := defaultCodec.Decode([]byte("test"), struct{}{})
		assert.Error(t, err)
		assert.Equal(t, "data without schema cannot be decoded", err.Error())
	})

	t.Run("avro binary encoding and decoding success", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaAvro, pubsub.EncodingBinary, getTestAvroSchemaDefinition(t))

		in := avro.SimpleRecord{
			StringField:  "test",
			FloatField:   12.34,
			BooleanField: true,
		}

		enc, err := defaultCodec.Encode(in)
		assert.NoError(t, err)

		out := avro.SimpleRecord{}

		err = defaultCodec.Decode(enc, &out)
		assert.NoError(t, err)

		assert.Equal(t, in, out)
	})

	t.Run("avro binary encoding failure", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaAvro, pubsub.EncodingBinary, getTestAvroSchemaDefinition(t))

		in := avro.InvalidSimpleRecord{
			StringField:  true,
			FloatField:   "test",
			BooleanField: 12.34,
		}

		_, err := defaultCodec.Encode(in)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot encode avro binary")
	})

	t.Run("avro binary decoding failure", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaAvro, pubsub.EncodingBinary, getTestAvroSchemaDefinition(t))

		out := avro.SimpleRecord{}

		err := defaultCodec.Decode([]byte("invalid"), &out)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot decode avro binary")
	})

	t.Run("avro json encoding and decoding success", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaAvro, pubsub.EncodingJSON, getTestAvroSchemaDefinition(t))

		in := avro.SimpleRecord{
			StringField:  "test",
			FloatField:   12.34,
			BooleanField: true,
		}

		enc, err := defaultCodec.Encode(in)
		assert.NoError(t, err)

		out := avro.SimpleRecord{}

		err = defaultCodec.Decode(enc, &out)
		assert.NoError(t, err)

		assert.Equal(t, in, out)
	})

	t.Run("avro json encoding failure", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaAvro, pubsub.EncodingJSON, getTestAvroSchemaDefinition(t))

		in := avro.InvalidSimpleRecord{
			StringField:  true,
			FloatField:   "test",
			BooleanField: 12.34,
		}

		_, err := defaultCodec.Encode(in)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot encode avro json")
	})

	t.Run("avro json decoding failure", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaAvro, pubsub.EncodingJSON, getTestAvroSchemaDefinition(t))

		out := avro.SimpleRecord{}

		err := defaultCodec.Decode([]byte("invalid"), &out)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot decode avro json")
	})

	t.Run("protobuf binary encoding and decoding success", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaProtocolBuffer, pubsub.EncodingBinary, getTestProtoSchemaDefinition(t))

		in := &proto.SimpleRecord{
			StringField:  "test",
			FloatField:   12.34,
			BooleanField: true,
		}

		enc, err := defaultCodec.Encode(in)
		assert.NoError(t, err)

		out := proto.SimpleRecord{}

		err = defaultCodec.Decode(enc, &out)
		assert.NoError(t, err)

		assert.Equal(t, in.StringField, out.StringField)
		assert.Equal(t, in.FloatField, out.FloatField)
		assert.Equal(t, in.BooleanField, out.BooleanField)
	})

	t.Run("protobuf binary encoding failure", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaProtocolBuffer, pubsub.EncodingBinary, getTestProtoSchemaDefinition(t))

		_, err := defaultCodec.Encode(struct{}{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid proto message")
	})

	t.Run("protobuf binary decoding failure", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaProtocolBuffer, pubsub.EncodingBinary, getTestProtoSchemaDefinition(t))

		out := proto.SimpleRecord{}

		err := defaultCodec.Decode([]byte("invalid"), &out)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot decode proto binary")
	})

	t.Run("protobuf json encoding and decoding success", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaProtocolBuffer, pubsub.EncodingJSON, getTestProtoSchemaDefinition(t))

		in := &proto.SimpleRecord{
			StringField:  "test",
			FloatField:   12.34,
			BooleanField: true,
		}

		enc, err := defaultCodec.Encode(in)
		assert.NoError(t, err)

		out := proto.SimpleRecord{}

		err = defaultCodec.Decode(enc, &out)
		assert.NoError(t, err)

		assert.Equal(t, in.StringField, out.StringField)
		assert.Equal(t, in.FloatField, out.FloatField)
		assert.Equal(t, in.BooleanField, out.BooleanField)
	})

	t.Run("protobuf json encoding failure", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaProtocolBuffer, pubsub.EncodingJSON, getTestProtoSchemaDefinition(t))

		_, err := defaultCodec.Encode(struct{}{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid proto message")
	})

	t.Run("protobuf json decoding failure", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaProtocolBuffer, pubsub.EncodingJSON, getTestProtoSchemaDefinition(t))

		out := proto.SimpleRecord{}

		err := defaultCodec.Decode([]byte("invalid"), &out)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot decode proto json")
	})

	t.Run("invalid schema type encoding failure", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaType(99), pubsub.EncodingJSON, "")

		_, err := defaultCodec.Encode([]byte("invalid"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid schema type")
	})

	t.Run("avro invalid encoding failure", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaAvro, pubsub.SchemaEncoding(99), "")

		_, err := defaultCodec.Encode([]byte("invalid"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid avro encoding")
	})

	t.Run("proto invalid encoding failure", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaProtocolBuffer, pubsub.SchemaEncoding(99), "")

		_, err := defaultCodec.Encode([]byte("invalid"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid proto encoding")
	})

	t.Run("invalid schema type decoding failure", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaType(99), pubsub.EncodingJSON, "")

		err := defaultCodec.Decode([]byte("invalid"), struct{}{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid schema type")
	})

	t.Run("avro invalid decoding failure", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaAvro, pubsub.SchemaEncoding(99), "")

		err := defaultCodec.Decode([]byte("invalid"), struct{}{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid avro encoding")
	})

	t.Run("proto invalid decoding failure", func(t *testing.T) {
		t.Parallel()

		defaultCodec := codec.NewDefaultCodec(pubsub.SchemaProtocolBuffer, pubsub.SchemaEncoding(99), "")

		err := defaultCodec.Decode([]byte("invalid"), struct{}{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid proto encoding")
	})
}

func getTestAvroSchemaDefinition(tb testing.TB) string {
	tb.Helper()

	data, err := os.ReadFile("../testdata/avro/simple.avsc")
	assert.NoError(tb, err)

	return string(data)
}

func getTestProtoSchemaDefinition(tb testing.TB) string {
	tb.Helper()

	data, err := os.ReadFile("../testdata/proto/simple.proto")
	assert.NoError(tb, err)

	return string(data)
}
