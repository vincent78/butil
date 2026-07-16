package logger4

import (
	"github.com/vincent78/butil/bus/core/logger"
	"go.uber.org/zap/zapcore"
)

type Option = logger.Option
type FileOption = logger.FileOption

func WithLevel(level string) Option {
	return logger.WithLevel(level)
}

func WithFormat(format string) Option {
	return logger.WithFormat(format)
}

func WithCaller(disableCaller bool, skip int) Option {
	return logger.WithCaller(disableCaller, skip)
}

func WithDisableCaller(disableCaller bool) Option {
	return logger.WithCaller(disableCaller, 1)
}

func WithStacktraceLevel(level string) Option {
	return logger.WithStacktraceLevel(level)
}

func WithHooks(hooks ...func(zapcore.Entry) error) Option {
	return logger.WithHooks(hooks...)
}

func WithSave(isSave bool, opts ...FileOption) Option {
	return logger.WithSave(isSave, opts...)
}

func WithFileName(filename string) FileOption {
	return logger.WithFileName(filename)
}

func WithFileMaxSize(maxSize int) FileOption {
	return logger.WithFileMaxSize(maxSize)
}

func WithFileMaxBackups(maxBackups int) FileOption {
	return logger.WithFileMaxBackups(maxBackups)
}

func WithFileMaxAge(maxAge int) FileOption {
	return logger.WithFileMaxAge(maxAge)
}

func WithFileIsCompression(isCompression bool) FileOption {
	return logger.WithFileIsCompression(isCompression)
}

func WithLocalTime(isLocalTime bool) FileOption {
	return logger.WithLocalTime(isLocalTime)
}
