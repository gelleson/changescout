package transform

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
)

func expectedHash(input []any) string {
	hasher := sha256.New()
	for _, elem := range input {
		hasher.Write([]byte(fmt.Sprintf("%v", elem)))
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func TestHashSlice(t *testing.T) {
	tests := []struct {
		name   string
		input  []any
		output string
	}{
		{
			name:   "empty slice",
			input:  []any{},
			output: expectedHash([]any{}),
		},
		{
			name:   "integer slice",
			input:  []any{1, 2, 3},
			output: expectedHash([]any{1, 2, 3}),
		},
		{
			name:   "string slice",
			input:  []any{"a", "b", "c"},
			output: expectedHash([]any{"a", "b", "c"}),
		},
		{
			name:   "mixed slice",
			input:  []any{1, "a", 3.5},
			output: expectedHash([]any{1, "a", 3.5}),
		},
		{
			name:   "duplicate elements",
			input:  []any{"a", "a", "b"},
			output: expectedHash([]any{"a", "a", "b"}),
		},
		{
			name:   "single element",
			input:  []any{42},
			output: expectedHash([]any{42}),
		},
		{
			name:   "long slice",
			input:  make([]any, 1000),
			output: expectedHash(make([]any, 1000)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashSlice(tt.input); got != tt.output {
				t.Errorf("HashSlice() = %v, want %v", got, tt.output)
			}
		})
	}
}
