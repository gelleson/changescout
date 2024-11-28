package logger

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}

func strLevelToZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func SetLevel(level string) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(strLevelToZapLevel(level))
	logger, _ = cfg.Build()
}

func L(name ...string) *zap.Logger {
	if len(name) > 0 {
		return logger.Named(name[0])
	}
	return logger
}

type loggerKey struct {
}

func WithLogger(_logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(
				c.Request().WithContext(
					context.WithValue(c.Request().Context(), loggerKey{}, _logger),
				),
			)
			return next(c)
		}
	}
}

func FromContext(c context.Context) *zap.Logger {
	_logger, ok := c.Value(loggerKey{}).(*zap.Logger)
	if ok {
		return _logger
	}
	return logger
}
