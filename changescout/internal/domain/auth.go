package domain

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

// AuthClaims represents the claims within a JWT.
type AuthClaims struct {
	jwt.RegisteredClaims `json:"inline"`
	ID                   uuid.UUID `json:"id"`
	Email                string    `json:"email"`
	ExpiresAt            time.Time `json:"exp"`
}

// AuthError represents an error during authentication.
type AuthError struct {
	Err error
}

func (e *AuthError) Error() string {
	return e.Err.Error()
}

func (e *AuthError) Unwrap() error {
	return e.Err
}

// IsUnauthenticated checks if the error is an authentication failure.
func IsUnauthenticated(err error) bool {
	var authErr *AuthError
	return errors.As(err, &authErr)
}

var _ jwt.Claims = (*AuthClaims)(nil)
