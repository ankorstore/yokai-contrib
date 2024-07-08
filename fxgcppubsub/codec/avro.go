package codec

import (
	"encoding/json"
	"fmt"

	"github.com/hamba/avro/v2"
	"github.com/linkedin/goavro/v2"
)

var (
	_ Codec = (*AvroBinaryCodec)(nil)
	_ Codec = (*AvroJsonCodec)(nil)
)

// AvroBinaryCodec is a Codec implementation for encoding and decoding with avro schema in binary format.
type AvroBinaryCodec struct {
	api    avro.API
	schema avro.Schema
}

// NewAvroBinaryCodec returns a new AvroBinaryCodec instance.
func NewAvroBinaryCodec(schemaDefinition string) (*AvroBinaryCodec, error) {
	api := avro.Config{}.Freeze()

	schema, err := avro.Parse(schemaDefinition)
	if err != nil {
		return nil, fmt.Errorf("cannot parse avro schema: %w", err)
	}

	return &AvroBinaryCodec{
		api:    api,
		schema: schema,
	}, nil
}

// Encode encodes in avro binary format the provided input.
func (c *AvroBinaryCodec) Encode(in any) ([]byte, error) {
	out, err := c.api.Marshal(c.schema, in)
	if err != nil {
		return nil, fmt.Errorf("cannot encode avro binary: %w", err)
	}

	return out, nil
}

// Decode decodes from avro binary format the provided input.
func (c *AvroBinaryCodec) Decode(enc []byte, out any) error {
	err := c.api.Unmarshal(c.schema, enc, out)
	if err != nil {
		return fmt.Errorf("cannot decode avro binary: %w", err)
	}

	return nil
}

// AvroJsonCodec is a Codec implementation for encoding and decoding with avro schema in json format.
type AvroJsonCodec struct {
	codec *goavro.Codec
}

// NewAvroJsonCodec returns a new AvroJsonCodec instance.
func NewAvroJsonCodec(schemaDefinition string) (*AvroJsonCodec, error) {
	codec, err := goavro.NewCodec(schemaDefinition)
	if err != nil {
		return nil, fmt.Errorf("cannot parse avro schema: %w", err)
	}

	return &AvroJsonCodec{codec: codec}, nil
}

// Encode encodes in avro json format.
func (c *AvroJsonCodec) Encode(in any) ([]byte, error) {
	inMap, err := c.convertStructIntoMap(in)
	if err != nil {
		return nil, fmt.Errorf("cannot convert struct into map: %w", err)
	}

	out, err := c.codec.TextualFromNative(nil, inMap)
	if err != nil {
		return nil, fmt.Errorf("cannot encode avro json: %w", err)
	}

	return out, nil
}

// Decode decodes from avro json format.
func (c *AvroJsonCodec) Decode(enc []byte, out any) error {
	data, _, err := c.codec.NativeFromTextual(enc)
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

func (c *AvroJsonCodec) convertStructIntoMap(in any) (map[string]interface{}, error) {
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

func (c *AvroJsonCodec) convertMapIntoStruct(in map[string]interface{}, out any) error {
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
