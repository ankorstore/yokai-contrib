package codec_test

import (
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/avro"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/testdata/proto"
	"github.com/stretchr/testify/assert"
)

func TestDefaultCodecFactory(t *testing.T) {
	t.Parallel()

	avroSchemaDefinition := avro.GetTestAvroSchemaDefinition(t)
	protoSchemaDefinition := proto.GetTestProtoSchemaDefinition(t)

	t.Run("construction", func(t *testing.T) {
		t.Parallel()

		defaultFactory := codec.NewDefaultCodecFactory()

		assert.IsType(t, &codec.DefaultCodecFactory{}, defaultFactory)
		assert.Implements(t, (*codec.CodecFactory)(nil), defaultFactory)
	})

	t.Run("raw codec creation", func(t *testing.T) {
		t.Parallel()

		cod, err := codec.NewDefaultCodecFactory().Create(
			pubsub.SchemaTypeUnspecified,
			pubsub.EncodingUnspecified,
			"",
		)
		assert.NoError(t, err)
		assert.IsType(t, &codec.RawCodec{}, cod)
		assert.Implements(t, (*codec.Codec)(nil), cod)
	})

	t.Run("avro binary codec creation", func(t *testing.T) {
		t.Parallel()

		cod, err := codec.NewDefaultCodecFactory().Create(
			pubsub.SchemaAvro,
			pubsub.EncodingBinary,
			avroSchemaDefinition,
		)
		assert.NoError(t, err)
		assert.IsType(t, &codec.AvroBinaryCodec{}, cod)
		assert.Implements(t, (*codec.Codec)(nil), cod)
	})

	t.Run("avro json codec creation", func(t *testing.T) {
		t.Parallel()

		cod, err := codec.NewDefaultCodecFactory().Create(
			pubsub.SchemaAvro,
			pubsub.EncodingJSON,
			avroSchemaDefinition,
		)
		assert.NoError(t, err)
		assert.IsType(t, &codec.AvroJsonCodec{}, cod)
		assert.Implements(t, (*codec.Codec)(nil), cod)
	})

	t.Run("proto binary codec creation", func(t *testing.T) {
		t.Parallel()

		cod, err := codec.NewDefaultCodecFactory().Create(
			pubsub.SchemaProtocolBuffer,
			pubsub.EncodingBinary,
			protoSchemaDefinition,
		)
		assert.NoError(t, err)
		assert.IsType(t, &codec.ProtoBinaryCodec{}, cod)
		assert.Implements(t, (*codec.Codec)(nil), cod)
	})

	t.Run("proto json codec creation", func(t *testing.T) {
		t.Parallel()

		cod, err := codec.NewDefaultCodecFactory().Create(
			pubsub.SchemaProtocolBuffer,
			pubsub.EncodingJSON,
			protoSchemaDefinition,
		)
		assert.NoError(t, err)
		assert.IsType(t, &codec.ProtoJsonCodec{}, cod)
		assert.Implements(t, (*codec.Codec)(nil), cod)
	})

	t.Run("invalid schema type codec creation", func(t *testing.T) {
		t.Parallel()

		cod, err := codec.NewDefaultCodecFactory().Create(
			pubsub.SchemaType(99),
			pubsub.EncodingJSON,
			"",
		)
		assert.Nil(t, cod)
		assert.Error(t, err)
		assert.Equal(t, "invalid schema type", err.Error())
	})

	t.Run("invalid avro encoding codec creation", func(t *testing.T) {
		t.Parallel()

		cod, err := codec.NewDefaultCodecFactory().Create(
			pubsub.SchemaAvro,
			pubsub.SchemaEncoding(99),
			"",
		)
		assert.Nil(t, cod)
		assert.Error(t, err)
		assert.Equal(t, "invalid avro encoding", err.Error())
	})

	t.Run("invalid proto encoding codec creation", func(t *testing.T) {
		t.Parallel()

		cod, err := codec.NewDefaultCodecFactory().Create(
			pubsub.SchemaProtocolBuffer,
			pubsub.SchemaEncoding(99),
			"",
		)
		assert.Nil(t, cod)
		assert.Error(t, err)
		assert.Equal(t, "invalid proto encoding", err.Error())
	})
}
