package message

import (
	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
)

// Message represents a pub/sub message with an associated codec.Codec.
type Message struct {
	codec   codec.Codec
	message *pubsub.Message
}

// NewMessage returns a new Message instance.
func NewMessage(codec codec.Codec, message *pubsub.Message) *Message {
	return &Message{
		codec:   codec,
		message: message,
	}
}

// Codec returns the associated codec.Codec.
func (m *Message) Codec() codec.Codec {
	return m.codec
}

// BaseMessage returns the base pubsub.Message.
func (m *Message) BaseMessage() *pubsub.Message {
	return m.message
}

// Decode decodes the message content into the provided parameter.
func (m *Message) Decode(out any) error {
	return m.codec.Decode(m.message.Data, out)
}

// ID returns the base message id.
func (m *Message) ID() string {
	return m.message.ID
}

// Data returns the base message data.
func (m *Message) Data() []byte {
	return m.message.Data
}

// Attributes returns the base message attributes.
func (m *Message) Attributes() map[string]string {
	return m.message.Attributes
}

// Ack indicates the successful message processing.
// Calls to Ack or Nack have no effect after the first call.
func (m *Message) Ack() {
	m.message.Ack()
}

// Nack indicates that the client will not or cannot process the message.
// Calls to Ack or Nack have no effect after the first call.
func (m *Message) Nack() {
	m.message.Nack()
}
