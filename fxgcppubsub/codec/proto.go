package codec

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	_ Codec = (*ProtoBinaryCodec)(nil)
	_ Codec = (*ProtoJsonCodec)(nil)
)

// ProtoBinaryCodec is a Codec implementation for encoding and decoding with protobuf schema in binary format.
type ProtoBinaryCodec struct {
	schemaDefinition string
}

// NewProtoBinaryCodec returns a new AvroBinaryCodec instance.
func NewProtoBinaryCodec(schemaDefinition string) *ProtoBinaryCodec {
	return &ProtoBinaryCodec{
		schemaDefinition: schemaDefinition,
	}
}

// Encode encodes in protobuf binary format the provided input.
func (c *ProtoBinaryCodec) Encode(in any) ([]byte, error) {
	protoIn, ok := in.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("invalid proto message")
	}

	out, err := proto.Marshal(protoIn)
	if err != nil {
		return nil, fmt.Errorf("cannot encode proto binary: %w", err)
	}

	return out, nil
}

// Decode decodes from protobuf binary format the provided input.
func (c *ProtoBinaryCodec) Decode(enc []byte, out any) error {
	protoOut, ok := out.(proto.Message)
	if !ok {
		return fmt.Errorf("invalid proto message")
	}

	err := proto.Unmarshal(enc, protoOut)
	if err != nil {
		return fmt.Errorf("cannot decode proto binary: %w", err)
	}

	return nil
}

// ProtoJsonCodec is a Codec implementation for encoding and decoding with protobuf schema in json format.
type ProtoJsonCodec struct {
	schemaDefinition string
}

// NewProtoJsonCodec returns a new AvroBinaryCodec instance.
func NewProtoJsonCodec(schemaDefinition string) *ProtoJsonCodec {
	return &ProtoJsonCodec{
		schemaDefinition: schemaDefinition,
	}
}

// Encode encodes in protobuf json format.
func (c *ProtoJsonCodec) Encode(in any) ([]byte, error) {
	protoIn, ok := in.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("invalid proto message")
	}

	out, err := protojson.Marshal(protoIn)
	if err != nil {
		return nil, fmt.Errorf("cannot encode proto json: %w", err)
	}

	return out, nil
}

// Decode decodes from protobuf json format.
func (c *ProtoJsonCodec) Decode(enc []byte, out any) error {
	protoOut, ok := out.(proto.Message)
	if !ok {
		return fmt.Errorf("invalid proto message")
	}

	err := protojson.Unmarshal(enc, protoOut)
	if err != nil {
		return fmt.Errorf("cannot decode proto json: %w", err)
	}

	return nil
}
