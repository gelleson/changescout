package auth

import (
	"github.com/gelleson/changescout/changescout/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// CheckPassword compares a hashed password with a plain text password using bcrypt.
// Returns an error if the passwords do not match, otherwise returns nil.
func CheckPassword(hashedPassword, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return domain.ErrInvalidCredentials
	}
	return nil
}

// HashPassword hashes a plain text password using bcrypt.
// Returns the hashed password as a string.
func HashPassword(plainPassword string, cost int) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
