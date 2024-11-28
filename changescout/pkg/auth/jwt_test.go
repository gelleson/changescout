package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// CustomClaims defines custom claims structure for test purpose
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func TestGenerateJWTWithClaims(t *testing.T) {
	// Test secret key
	secret := []byte("mysecret")

	// Define the claims
	claims := CustomClaims{
		UserID: "12345",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	// Generate the token
	tokenString, err := GenerateJWTWithClaims(claims, secret)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parse the token to verify
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	assert.NoError(t, err)

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		assert.Equal(t, "12345", claims.UserID)
	} else {
		t.Fatalf("Token is invalid or claims are incorrect")
	}
}
