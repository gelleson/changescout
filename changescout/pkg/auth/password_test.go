package auth

import (
	"testing"

	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// TestCheckPassword tests the CheckPassword function.
func TestCheckPassword(t *testing.T) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name           string
		hashedPassword string
		plainPassword  string
		expectedErr    error
	}{
		{"Valid password", string(hashedPassword), "password123", nil},
		{"Invalid password", string(hashedPassword), "wrongpassword", domain.ErrInvalidCredentials},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPassword(tt.hashedPassword, tt.plainPassword)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

// TestHashPassword tests the HashPassword function.
func TestHashPassword(t *testing.T) {
	plainPassword := "password123"
	cost := bcrypt.DefaultCost

	hashedPassword, err := HashPassword(plainPassword, cost)
	assert.NoError(t, err)

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	assert.NoError(t, err, "Hashed password should match the plain password")
}
