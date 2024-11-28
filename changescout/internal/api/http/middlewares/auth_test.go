package middlewares

import (
	"github.com/gelleson/changescout/changescout/pkg/auth"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/pkg/contexts"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type UserContextKey struct{}

// Custom claims for testing purposes.
type testClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

func (c testClaims) Valid() error {
	return nil
}

// Test JWTAuth middleware.
func TestJWTAuth(t *testing.T) {
	e := echo.New()

	secret := []byte("secret")
	email := "test@example.com"
	user := &domain.AuthClaims{
		ID:    uuid.New(),
		Email: email,
	}

	cfg := JWTAuthConfig{
		Secret: secret,
	}

	middleware := JWTAuth(cfg)

	t.Run("Valid token", func(t *testing.T) {
		claims := domain.AuthClaims{
			ID:        user.ID,
			Email:     user.Email,
			ExpiresAt: time.Now().Add(time.Hour),
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:   "changescout",
				Subject:  "changescout",
				Audience: []string{"changescout"},
			},
		}
		token, err := auth.GenerateJWTWithClaims(claims, secret)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Cookie", "accessToken="+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		err = handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		extractedUser := c.Request().Context().Value(contexts.UserContextKey{}).(*domain.AuthClaims)
		assert.Equal(t, user.ID, extractedUser.ID)
	})

	t.Run("Invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer invalidtoken")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Ensure the user context is not set
		extractedUser := c.Request().Context().Value(contexts.UserContextKey{})
		assert.Nil(t, extractedUser)
	})
}
