package messages

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string
	Value int
}

func TestEncode(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name:  "encode_struct",
			input: TestStruct{Name: "test", Value: 123},
			want:  `{"Name":"test","Value":123}`,
		},
		{
			name:  "encode_map",
			input: map[string]int{"key": 1, "another": 2},
			want:  `{"another":2,"key":1}`,
		},
		{
			name:  "encode_slice",
			input: []string{"a", "b", "c"},
			want:  `["a","b","c"]`,
		},
		{
			name:  "encode_nil",
			input: nil,
			want:  `null`,
		},
		{
			name:  "encode_string",
			input: "string value",
			want:  `"string value"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Encode(tt.input)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestEncode_InvalidInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	Encode(make(chan int)) // panic
}

func TestNewMessage(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{
			name:  "new_message_struct",
			input: TestStruct{Name: "msg", Value: 101},
		},
		{
			name:  "new_message_slice",
			input: []int{1, 2, 3, 4},
		},
		{
			name:  "new_message_string",
			input: "simple message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := NewMessage(tt.input)
			assert.NotNil(t, msg)
			assert.NotEmpty(t, msg.UUID)

			encodedInput, _ := json.Marshal(tt.input)
			assert.JSONEq(t, string(encodedInput), string(msg.Payload))
		})
	}
}
