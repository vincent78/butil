package logger

import (
	"context"

	"go.uber.org/zap/zapcore"
)

// Field type
type Field = zapcore.Field

type ILoggerMethod interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	DPanic(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Sync() error
	WithMetaCtx(ctx context.Context, keys ...string) ILoggerMethod
	WithFields(fields ...Field) ILoggerMethod
	WithMap(map[string]any) ILoggerMethod
	WithCallerSkip(skip int) ILoggerMethod
}
