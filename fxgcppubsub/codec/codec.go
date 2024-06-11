package codec

import (
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/hamba/avro/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Codec struct {
	schemaConfig   *pubsub.SchemaConfig
	schemaSettings *pubsub.SchemaSettings
}

func NewCodec(schemaConfig *pubsub.SchemaConfig, schemaSettings *pubsub.SchemaSettings) *Codec {
	return &Codec{
		schemaConfig:   schemaConfig,
		schemaSettings: schemaSettings,
	}
}

func (c *Codec) Encode(in any) ([]byte, error) {
	if c.schemaConfig == nil || c.schemaSettings == nil {
		inBytes, ok := in.([]byte)
		if !ok {
			return nil, fmt.Errorf("data without schema must be of type []byte")
		}

		return inBytes, nil
	}

	var out []byte
	var err error

	switch c.schemaConfig.Type {
	case pubsub.SchemaAvro:
		switch c.schemaSettings.Encoding {
		case pubsub.EncodingBinary:
			out, err = c.encodeAvroBinary(in)
		case pubsub.EncodingJSON:
			out, err = c.encodeAvroJSON(in)
		default:
			err = fmt.Errorf("invalid avro encoding")
		}
	case pubsub.SchemaProtocolBuffer:
		switch c.schemaSettings.Encoding {
		case pubsub.EncodingBinary:
			out, err = c.encodeProtoBinary(in)
		case pubsub.EncodingJSON:
			out, err = c.encodeProtoJSON(in)
		default:
			err = fmt.Errorf("invalid proto encoding")
		}
	default:
		err = fmt.Errorf("invalid schema type")
	}

	return out, err
}

func (c *Codec) Decode(enc []byte, out any) error {
	if c.schemaConfig == nil || c.schemaSettings == nil {
		return fmt.Errorf("no schema associated, nothing to decode, use message data instead")
	}

	var err error

	switch c.schemaConfig.Type {
	case pubsub.SchemaAvro:
		switch c.schemaSettings.Encoding {
		case pubsub.EncodingBinary:
			err = c.decodeAvroBinary(enc, out)
		case pubsub.EncodingJSON:
			err = c.decodeAvroJSON(enc, out)
		default:
			err = fmt.Errorf("invalid avro encoding")
		}
	case pubsub.SchemaProtocolBuffer:
		switch c.schemaSettings.Encoding {
		case pubsub.EncodingBinary:
			err = c.decodeProtoBinary(enc, out)
		case pubsub.EncodingJSON:
			err = c.decodeProtoJSON(enc, out)
		default:
			err = fmt.Errorf("invalid proto encoding")
		}
	default:
		err = fmt.Errorf("invalid schema type")
	}

	return err
}

func (c *Codec) encodeAvroBinary(in any) ([]byte, error) {
	avroSchema, err := avro.Parse(c.schemaConfig.Definition)
	if err != nil {
		return nil, fmt.Errorf("cannot parse avro schema: %w", err)
	}

	out, err := avro.Marshal(avroSchema, in)
	if err != nil {
		return nil, fmt.Errorf("cannot encode avro binary: %w", err)
	}

	return out, nil
}

func (c *Codec) decodeAvroBinary(enc []byte, out any) error {
	avroSchema, err := avro.Parse(c.schemaConfig.Definition)
	if err != nil {
		return fmt.Errorf("cannot parse avro schema: %w", err)
	}

	err = avro.Unmarshal(avroSchema, enc, out)
	if err != nil {
		return fmt.Errorf("cannot decode avro binary: %w", err)
	}

	return nil
}

func (c *Codec) encodeAvroJSON(in any) ([]byte, error) {
	out, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("cannot encode avro json: %w", err)
	}

	return out, nil
}

func (c *Codec) decodeAvroJSON(enc []byte, out any) error {
	err := json.Unmarshal(enc, out)
	if err != nil {
		return fmt.Errorf("cannot decode avro json: %w", err)
	}

	return nil
}

func (c *Codec) encodeProtoBinary(in any) ([]byte, error) {
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

func (c *Codec) decodeProtoBinary(enc []byte, out any) error {
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

func (c *Codec) encodeProtoJSON(in any) ([]byte, error) {
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

func (c *Codec) decodeProtoJSON(enc []byte, out any) error {
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
