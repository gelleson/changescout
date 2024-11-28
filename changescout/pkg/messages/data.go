package messages

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
)

type MessageWrapper[T any] struct {
	*message.Message
}

func (m *MessageWrapper[T]) Get() T {
	var v T
	_ = json.Unmarshal(m.Payload, &v)
	return v
}

func Wrap[T any](msg *message.Message) *MessageWrapper[T] {
	return &MessageWrapper[T]{
		Message: msg,
	}
}
