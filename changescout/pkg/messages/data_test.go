package messages

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"reflect"
	"testing"
)

func TestMessageWrapper_Get(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 int
	}

	tests := []struct {
		name    string
		payload []byte
		want    TestStruct
	}{
		{name: "ValidJSON", payload: []byte(`{"Field1":"value","Field2":123}`), want: TestStruct{"value", 123}},
		{name: "MissingFields", payload: []byte(`{"Field1":"value"}`), want: TestStruct{"value", 0}},
		{name: "ExtraFields", payload: []byte(`{"Field1":"value","Field2":123,"Field3":"extra"}`), want: TestStruct{"value", 123}},
		{name: "EmptyJSON", payload: []byte(`{}`), want: TestStruct{"", 0}},
		{name: "InvalidJSON", payload: []byte(`{"Field1":"value",`), want: TestStruct{"", 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &message.Message{Payload: tt.payload}
			wrapper := Wrap[TestStruct](msg)
			if got := wrapper.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MessageWrapper.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		name string
		msg  *message.Message
	}{
		{name: "SimpleWrap", msg: &message.Message{Payload: []byte(`{"Field1":"test"}`)}},
		{name: "EmptyPayload", msg: &message.Message{Payload: []byte(``)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Wrap[interface{}](tt.msg)
			if !reflect.DeepEqual(got.Message, tt.msg) {
				t.Errorf("Wrap() = %v, want %v", got.Message, tt.msg)
			}
		})
	}
}
