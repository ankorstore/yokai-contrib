package message

import (
	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/codec"
)

type Message struct {
	message *pubsub.Message
	codec   *codec.Codec
}

func NewMessage(message *pubsub.Message, codec *codec.Codec) *Message {
	return &Message{
		message: message,
		codec:   codec,
	}
}

func (m *Message) Base() *pubsub.Message {
	return m.message
}

func (m *Message) Decode(out any) error {
	return m.codec.Decode(m.message.Data, out)
}

func (m *Message) ID() string {
	return m.message.ID
}

func (m *Message) Data() []byte {
	return m.message.Data
}

func (m *Message) Attributes() map[string]string {
	return m.message.Attributes
}

func (m *Message) Ack() {
	m.message.Ack()
}

func (m *Message) Nack() {
	m.message.Nack()
}
