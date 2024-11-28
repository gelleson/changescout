package middlewares

import (
	"errors"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/pkg/contexts"
	"github.com/gelleson/changescout/changescout/internal/platform/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type JWTAuthConfig struct {
	Secret []byte
}

func JWTAuth(cfg JWTAuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log := logger.FromContext(c.Request().Context())
			authHeader, err := c.Cookie("accessToken")
			if err != nil {
				log.Error("JWT verification failed", zap.Error(err))
				return next(c)
			}

			usr, err := verifyJWT(authHeader.Value, cfg.Secret)
			if err != nil {
				log.Error("JWT verification failed", zap.Error(err))
				return next(c)
			}

			c.SetRequest(
				c.Request().WithContext(
					contexts.WithUserContext(c.Request().Context(), usr),
				),
			)

			return next(c)
		}
	}
}

func verifyJWT(token string, secret []byte) (*domain.AuthClaims, error) {
	var zero domain.AuthClaims
	t, err := jwt.ParseWithClaims(token, &zero, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return &zero, err
	}
	if !t.Valid {
		return &zero, errors.New("invalid token")
	}
	return t.Claims.(*domain.AuthClaims), nil
}
