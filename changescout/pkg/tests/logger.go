package tests

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestLogger(t *testing.T) *zap.Logger {
	return zaptest.NewLogger(t, zaptest.Level(zap.FatalLevel))
}
