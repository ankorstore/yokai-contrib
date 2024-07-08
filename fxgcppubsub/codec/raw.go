package codec

import "fmt"

var (
	_ Codec = (*RawCodec)(nil)
)

// RawCodec is a Codec implementation for encoding and decoding without specified schemas.
type RawCodec struct{}

// NewRawCodec returns a new RawCodec instance.
func NewRawCodec() *RawCodec {
	return &RawCodec{}
}

// Encode encodes in []byte.
func (c *RawCodec) Encode(in any) ([]byte, error) {
	return []byte(fmt.Sprintf("%s", in)), nil
}

// Decode returns a error.
func (c *RawCodec) Decode([]byte, any) error {
	return fmt.Errorf("data without schema cannot be decoded")
}
