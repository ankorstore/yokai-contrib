package codec

import "cloud.google.com/go/pubsub"

var _ CodecFactory = (*DefaultCodecFactory)(nil)

// CodecFactory is the interface for Codec factories.
type CodecFactory interface {
	Create(schemaType pubsub.SchemaType, schemaEncoding pubsub.SchemaEncoding, schemaDefinition string) Codec
}

// DefaultCodecFactory is the default CodecFactory implementation.
type DefaultCodecFactory struct{}

// NewDefaultCodecFactory returns a new DefaultCodecFactory instance.
func NewDefaultCodecFactory() *DefaultCodecFactory {
	return &DefaultCodecFactory{}
}

// Create creates a new Codec.
func (f *DefaultCodecFactory) Create(schemaType pubsub.SchemaType, schemaEncoding pubsub.SchemaEncoding, schemaDefinition string) Codec {
	return NewDefaultCodec(schemaType, schemaEncoding, schemaDefinition)
}
