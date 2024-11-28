package http

import (
	"context"
	"github.com/labstack/echo/v4"
)

type echoContextKey struct{}

func WithEchoContext(ctx context.Context, ectx echo.Context) context.Context {
	return context.WithValue(ctx, echoContextKey{}, ectx)
}

func EchoFromContext(ctx context.Context) (echo.Context, bool) {
	e, ok := ctx.Value(echoContextKey{}).(echo.Context)
	return e, ok
}

func WithCtx() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(c.Request().WithContext(
				WithEchoContext(c.Request().Context(), c),
			))
			return next(c)
		}
	}
}
