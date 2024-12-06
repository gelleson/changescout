package database

import (
	"errors"
	"testing"
)

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"IsErrEntityNotFound", ErrEntityNotFound, true},
		{"IsWrappedErrEntityNotFound", errors.New("wrapped: " + ErrEntityNotFound.Error()), false},
		{"IsNilError", nil, false},
		{"IsDifferentError", errors.New("different error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFound(tt.err); got != tt.want {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}
