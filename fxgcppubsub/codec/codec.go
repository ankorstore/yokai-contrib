package codec

import (
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/hamba/avro/v2"
	"github.com/linkedin/goavro/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var _ Codec = (*DefaultCodec)(nil)

// Codec is the interface for components in charge to handle raw, avro or protobuf encoding and decoding.
type Codec interface {
	Encode(in any) ([]byte, error)
	Decode(enc []byte, out any) error
}

// DefaultCodec is the default Codec implementation.
type DefaultCodec struct {
	schemaType       pubsub.SchemaType
	schemaEncoding   pubsub.SchemaEncoding
	schemaDefinition string
}

// NewDefaultCodec returns a new DefaultCodec instance.
func NewDefaultCodec(schemaType pubsub.SchemaType, schemaEncoding pubsub.SchemaEncoding, schemaDefinition string) *DefaultCodec {
	return &DefaultCodec{
		schemaType:       schemaType,
		schemaEncoding:   schemaEncoding,
		schemaDefinition: schemaDefinition,
	}
}

// Encode encodes an input into an avro or protobuf encoded slice of bytes.
func (c *DefaultCodec) Encode(in any) ([]byte, error) {
	switch c.schemaType {
	case pubsub.SchemaTypeUnspecified:
		return []byte(fmt.Sprintf("%s", in)), nil
	case pubsub.SchemaAvro:
		switch c.schemaEncoding {
		case pubsub.EncodingBinary:
			return c.encodeAvroBinary(in)
		case pubsub.EncodingJSON:
			return c.encodeAvroJSON(in)
		default:
			return nil, fmt.Errorf("invalid avro encoding")
		}
	case pubsub.SchemaProtocolBuffer:
		switch c.schemaEncoding {
		case pubsub.EncodingBinary:
			return c.encodeProtoBinary(in)
		case pubsub.EncodingJSON:
			return c.encodeProtoJSON(in)
		default:
			return nil, fmt.Errorf("invalid proto encoding")
		}
	default:
		return nil, fmt.Errorf("invalid schema type")
	}
}

// Decode decodes an avro or protobuf encoded slice of bytes into a provided output.
func (c *DefaultCodec) Decode(enc []byte, out any) error {
	switch c.schemaType {
	case pubsub.SchemaTypeUnspecified:
		return fmt.Errorf("data without schema cannot be decoded")
	case pubsub.SchemaAvro:
		switch c.schemaEncoding {
		case pubsub.EncodingBinary:
			return c.decodeAvroBinary(enc, out)
		case pubsub.EncodingJSON:
			return c.decodeAvroJSON(enc, out)
		default:
			return fmt.Errorf("invalid avro encoding")
		}
	case pubsub.SchemaProtocolBuffer:
		switch c.schemaEncoding {
		case pubsub.EncodingBinary:
			return c.decodeProtoBinary(enc, out)
		case pubsub.EncodingJSON:
			return c.decodeProtoJSON(enc, out)
		default:
			return fmt.Errorf("invalid proto encoding")
		}
	default:
		return fmt.Errorf("invalid schema type")
	}
}

func (c *DefaultCodec) encodeAvroBinary(in any) ([]byte, error) {
	avroSchema, err := avro.Parse(c.schemaDefinition)
	if err != nil {
		return nil, fmt.Errorf("cannot parse avro schema: %w", err)
	}

	out, err := avro.Marshal(avroSchema, in)
	if err != nil {
		return nil, fmt.Errorf("cannot encode avro binary: %w", err)
	}

	return out, nil
}

func (c *DefaultCodec) decodeAvroBinary(enc []byte, out any) error {
	avroSchema, err := avro.Parse(c.schemaDefinition)
	if err != nil {
		return fmt.Errorf("cannot parse avro schema: %w", err)
	}

	err = avro.Unmarshal(avroSchema, enc, out)
	if err != nil {
		return fmt.Errorf("cannot decode avro binary: %w", err)
	}

	return nil
}

func (c *DefaultCodec) encodeAvroJSON(in any) ([]byte, error) {
	avroSchema, err := goavro.NewCodec(c.schemaDefinition)
	if err != nil {
		return nil, fmt.Errorf("cannot parse avro schema: %w", err)
	}

	inMap, err := c.convertStructIntoMap(in)
	if err != nil {
		return nil, fmt.Errorf("cannot convert struct into map: %w", err)
	}

	out, err := avroSchema.TextualFromNative(nil, inMap)
	if err != nil {
		return nil, fmt.Errorf("cannot encode avro json: %w", err)
	}

	return out, nil
}

func (c *DefaultCodec) decodeAvroJSON(enc []byte, out any) error {
	avroSchema, err := goavro.NewCodec(c.schemaDefinition)
	if err != nil {
		return fmt.Errorf("cannot parse avro schema: %w", err)
	}

	data, _, err := avroSchema.NativeFromTextual(enc)
	if err != nil {
		return fmt.Errorf("cannot decode avro json: %w", err)
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("cannot convert avro json into map: %w", err)
	}

	err = c.convertMapIntoStruct(dataMap, out)
	if err != nil {
		return fmt.Errorf("cannot convert map into struct: %w", err)
	}

	return nil
}

func (c *DefaultCodec) encodeProtoBinary(in any) ([]byte, error) {
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

func (c *DefaultCodec) decodeProtoBinary(enc []byte, out any) error {
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

func (c *DefaultCodec) encodeProtoJSON(in any) ([]byte, error) {
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

func (c *DefaultCodec) decodeProtoJSON(enc []byte, out any) error {
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

func (c *DefaultCodec) convertStructIntoMap(in any) (map[string]interface{}, error) {
	var out map[string]interface{}

	jsonIn, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal json: %w", err)
	}

	err = json.Unmarshal(jsonIn, &out)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal json: %w", err)
	}

	return out, nil
}

func (c *DefaultCodec) convertMapIntoStruct(in map[string]interface{}, out any) error {
	jsonIn, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("cannot marshal json: %w", err)
	}

	err = json.Unmarshal(jsonIn, &out)
	if err != nil {
		return fmt.Errorf("cannot unmarshal json: %w", err)
	}

	return nil
}
