package messages

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

func Encode[T any](v T) []byte {
	d, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return d
}

func NewMessage[T any](v T) *message.Message {
	return message.NewMessage(
		uuid.New().String(),
		Encode(v),
	)
}
