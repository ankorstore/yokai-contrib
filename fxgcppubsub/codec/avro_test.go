package codec_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/avro"
	"github.com/stretchr/testify/assert"
)

func TestAvroBinaryCodec(t *testing.T) {
	t.Parallel()

	schemaDefinition := avro.GetTestAvroSchemaDefinition(t)

	t.Run("avro binary encoding and decoding success", func(t *testing.T) {
		t.Parallel()

		avroBinaryCodec, err := codec.NewAvroBinaryCodec(schemaDefinition)
		assert.NoError(t, err)

		in := avro.SimpleRecord{
			StringField:  "test",
			FloatField:   12.34,
			BooleanField: true,
		}

		enc, err := avroBinaryCodec.Encode(in)
		assert.NoError(t, err)

		out := avro.SimpleRecord{}

		err = avroBinaryCodec.Decode(enc, &out)
		assert.NoError(t, err)

		assert.Equal(t, in, out)
	})

	t.Run("avro binary encoding failure", func(t *testing.T) {
		t.Parallel()

		avroBinaryCodec, err := codec.NewAvroBinaryCodec(schemaDefinition)
		assert.NoError(t, err)

		in := avro.InvalidSimpleRecord{
			StringField:  true,
			FloatField:   "test",
			BooleanField: 12.34,
		}

		_, err = avroBinaryCodec.Encode(in)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot encode avro binary")
	})

	t.Run("avro binary decoding failure", func(t *testing.T) {
		t.Parallel()

		avroBinaryCodec, err := codec.NewAvroBinaryCodec(schemaDefinition)
		assert.NoError(t, err)

		out := avro.SimpleRecord{}

		err = avroBinaryCodec.Decode([]byte("invalid"), &out)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot decode avro binary")
	})

	t.Run("avro proto invalid schema", func(t *testing.T) {
		t.Parallel()

		avroBinaryCodec, err := codec.NewAvroBinaryCodec("invalid")
		assert.Nil(t, avroBinaryCodec)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot parse avro schema")
	})
}

func TestAvroJsonCodec(t *testing.T) {
	t.Parallel()

	schemaDefinition := avro.GetTestAvroSchemaDefinition(t)

	t.Run("avro json encoding and decoding success", func(t *testing.T) {
		t.Parallel()

		avroJsonCodec, err := codec.NewAvroJsonCodec(schemaDefinition)
		assert.NoError(t, err)

		in := avro.SimpleRecord{
			StringField:  "test",
			FloatField:   12.34,
			BooleanField: true,
		}

		enc, err := avroJsonCodec.Encode(in)
		assert.NoError(t, err)

		out := avro.SimpleRecord{}

		err = avroJsonCodec.Decode(enc, &out)
		assert.NoError(t, err)

		assert.Equal(t, in, out)
	})

	t.Run("avro json encoding failure", func(t *testing.T) {
		t.Parallel()

		avroJsonCodec, err := codec.NewAvroJsonCodec(schemaDefinition)
		assert.NoError(t, err)

		in := avro.InvalidSimpleRecord{
			StringField:  true,
			FloatField:   "test",
			BooleanField: 12.34,
		}

		_, err = avroJsonCodec.Encode(in)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot encode avro json")
	})

	t.Run("avro json decoding failure", func(t *testing.T) {
		t.Parallel()

		avroJsonCodec, err := codec.NewAvroJsonCodec(schemaDefinition)
		assert.NoError(t, err)

		out := avro.SimpleRecord{}

		err = avroJsonCodec.Decode([]byte("invalid"), &out)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot decode avro json")
	})

	t.Run("avro json invalid schema", func(t *testing.T) {
		t.Parallel()

		avroJsonCodec, err := codec.NewAvroJsonCodec("invalid")
		assert.Nil(t, avroJsonCodec)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot parse avro schema")
	})
}
