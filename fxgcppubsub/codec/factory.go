package codec

import (
	"fmt"

	"cloud.google.com/go/pubsub"
)

var _ CodecFactory = (*DefaultCodecFactory)(nil)

// Codec is the interface for codecs in charge to handle encoding and decoding.
type Codec interface {
	Encode(in any) ([]byte, error)
	Decode(enc []byte, out any) error
}

// CodecFactory is the interface for Codec factories.
type CodecFactory interface {
	Create(schemaType pubsub.SchemaType, schemaEncoding pubsub.SchemaEncoding, schemaDefinition string) (Codec, error)
}

// DefaultCodecFactory is the default CodecFactory implementation.
type DefaultCodecFactory struct{}

// NewDefaultCodecFactory returns a new DefaultCodecFactory instance.
func NewDefaultCodecFactory() *DefaultCodecFactory {
	return &DefaultCodecFactory{}
}

// Create creates a new Codec for given schema type, encoding and definition.
//
//nolint:cyclop,exhaustive
func (f *DefaultCodecFactory) Create(schemaType pubsub.SchemaType, schemaEncoding pubsub.SchemaEncoding, schemaDefinition string) (Codec, error) {
	switch schemaType {
	case pubsub.SchemaTypeUnspecified:
		return NewRawCodec(), nil
	case pubsub.SchemaAvro:
		switch schemaEncoding {
		case pubsub.EncodingBinary:
			return NewAvroBinaryCodec(schemaDefinition)
		case pubsub.EncodingJSON:
			return NewAvroJsonCodec(schemaDefinition)
		default:
			return nil, fmt.Errorf("invalid avro encoding")
		}
	case pubsub.SchemaProtocolBuffer:
		switch schemaEncoding {
		case pubsub.EncodingBinary:
			return NewProtoBinaryCodec(), nil
		case pubsub.EncodingJSON:
			return NewProtoJsonCodec(), nil
		default:
			return nil, fmt.Errorf("invalid proto encoding")
		}
	default:
		return nil, fmt.Errorf("invalid schema type")
	}
}
