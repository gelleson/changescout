package domain

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestIsUnauthenticated(t *testing.T) {
	// Define table driven test cases
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"NoError", nil, false},
		{"DifferentError", errors.New("different error"), false},
		{"AuthError", &AuthError{errors.New("auth error")}, true},
	}

	// Iterate over each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsUnauthenticated(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestAuthError_Unwrap(t *testing.T) {
	originalErr := errors.New("some error")
	authErr := &AuthError{Err: originalErr}

	unwrappedErr := authErr.Unwrap()

	if unwrappedErr != originalErr {
		t.Errorf("expected %v, got %v", originalErr, unwrappedErr)
	}
}

func TestAuthError_Error(t *testing.T) {
	originalErr := errors.New("some error")
	authErr := &AuthError{Err: originalErr}

	if authErr.Error() != originalErr.Error() {
		t.Errorf("expected %v, got %v", originalErr.Error(), authErr.Error())
	}
}

// Additional test to demonstrate usage of AuthClaims struct
func TestAuthClaims(t *testing.T) {
	id := uuid.New()
	email := "test@example.com"
	expiresAt := time.Now().Add(1 * time.Hour)
	claims := AuthClaims{
		ID:        id,
		Email:     email,
		ExpiresAt: expiresAt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	if claims.ID != id {
		t.Errorf("expected %v, got %v", id, claims.ID)
	}
	if claims.Email != email {
		t.Errorf("expected %s, got %s", email, claims.Email)
	}
	if !claims.ExpiresAt.Equal(expiresAt) {
		t.Errorf("expected %v, got %v", expiresAt, claims.ExpiresAt)
	}
}
