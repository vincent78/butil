package logger4

import (
	"context"
	"strings"
)

type LoggerMethod interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	DPanic(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	WithMetaCtx(ctx context.Context, keys ...string)
	WithFields(fields ...Field)
	WithMap(map[string]any)
}

// Debug level information
func Debug(msg string, fields ...Field) {
	getLogger().Debug(msg, fields...)
}

// Info level information
func Info(msg string, fields ...Field) {
	getLogger().Info(msg, fields...)
}

// Warn level information
func Warn(msg string, fields ...Field) {
	getLogger().Warn(msg, fields...)
}

// Error level information
func Error(msg string, fields ...Field) {
	getLogger().Error(msg, fields...)
}

// Panic level information
func Panic(msg string, fields ...Field) {
	getLogger().Panic(msg, fields...)
}

// Fatal level information
func Fatal(msg string, fields ...Field) {
	getLogger().Fatal(msg, fields...)
}

// Debugf format level information
func Debugf(format string, a ...any) {
	getSugaredLogger().Debugf(format, a...)
}

// Infof format level information
func Infof(format string, a ...any) {
	getSugaredLogger().Infof(format, a...)
}

// Warnf format level information
func Warnf(format string, a ...any) {
	getSugaredLogger().Warnf(format, a...)
}

// Errorf format level information
func Errorf(format string, a ...any) {
	getSugaredLogger().Errorf(format, a...)
}

// Fatalf format level information
func Fatalf(format string, a ...any) {
	getSugaredLogger().Fatalf(format, a...)
}

// Sync flushing any buffered log entries, applications should take care to call Sync before exiting.
func Sync() error {
	_ = getSugaredLogger().Sync()
	err := getLogger().Sync()
	if err != nil && !strings.Contains(err.Error(), "/dev/stdout") {
		return err
	}
	return nil
}

// WithFields carrying field information
func WithFields(fields ...Field) *NormalLogger {
	//return GetWithSkip(0).With(fields...)
	return Get().WithFields(fields...)
}
