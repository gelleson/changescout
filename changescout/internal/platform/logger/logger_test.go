package logger

import (
	"context"
	"github.com/labstack/echo/v4"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
)

func TestStrLevelToZapLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  zapcore.Level
	}{
		{"DebugLevel", "debug", zapcore.DebugLevel},
		{"InfoLevel", "info", zapcore.InfoLevel},
		{"WarnLevel", "warn", zapcore.WarnLevel},
		{"ErrorLevel", "error", zapcore.ErrorLevel},
		{"DPanicLevel", "dpanic", zapcore.DPanicLevel},
		{"PanicLevel", "panic", zapcore.PanicLevel},
		{"FatalLevel", "fatal", zapcore.FatalLevel},
		{"DefaultInfoLevel", "unknown", zapcore.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strLevelToZapLevel(tt.level); got != tt.want {
				t.Errorf("strLevelToZapLevel(%v) = %v, want %v", tt.level, got, tt.want)
			}
		})
	}
}

func TestSetLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  zapcore.Level
	}{
		{"SetDebugLevel", "debug", zapcore.DebugLevel},
		{"SetInfoLevel", "info", zapcore.InfoLevel},
		{"SetWarnLevel", "warn", zapcore.WarnLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLevel(tt.level)
			got := logger.Core().Enabled(tt.want)
			if !got {
				t.Errorf("logger level %v is not enabled, expected to be enabled", tt.want)
			}
		})
	}
}

func TestL(t *testing.T) {
	t.Run("DefaultLogger", func(t *testing.T) {
		if L() != logger {
			t.Errorf("L() != logger")
		}
	})

	t.Run("NamedLogger", func(t *testing.T) {
		namedLogger := L("test")
		if namedLogger == nil || namedLogger.Core() != logger.Core() {
			t.Errorf("Named logger does not match the default logger")
		}
	})
}

func TestWithLogger(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	myLogger := zap.NewExample()
	h := WithLogger(myLogger)(func(c echo.Context) error {
		loggerFromContext := FromContext(c.Request().Context())
		if loggerFromContext != myLogger {
			t.Errorf("Logger from context does not match myLogger")
		}
		return nil
	})

	h(c)
}

func TestFromContext(t *testing.T) {
	t.Run("ExistingLogger", func(t *testing.T) {
		myLogger := zap.NewExample()
		ctx := context.WithValue(context.Background(), loggerKey{}, myLogger)
		loggerFromContext := FromContext(ctx)
		if loggerFromContext != myLogger {
			t.Errorf("Logger from context does not match myLogger")
		}
	})

	t.Run("DefaultLogger", func(t *testing.T) {
		loggerFromContext := FromContext(context.Background())
		if loggerFromContext != logger {
			t.Errorf("Logger from context does not match default logger")
		}
	})
}
