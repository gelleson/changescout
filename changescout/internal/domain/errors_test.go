package domain

import (
	"errors"
	"testing"
)

func TestIsErrCheckNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "ErrCheckNotFound",
			err:      ErrCheckNotFound,
			expected: true,
		},
		{
			name:     "WrappedErrCheckNotFound",
			err:      errors.New("wrapped: check not found"),
			expected: false,
		},
		{
			name:     "DifferentError",
			err:      errors.New("different error"),
			expected: false,
		},
		{
			name:     "NoError",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsErrCheckNotFound(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
