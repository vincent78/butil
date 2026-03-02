package logger4

import (
	"strings"

	"github.com/vincent78/butil/bus/core/logger"
)

/*************************************************************************************************
 *
 * default的外露方法
 *
 **************************************************************************************************/

// Debug level information
func Debug(msg string, fields ...logger.Field) {
	Get().Debug(msg, fields...)
}

// Info level information
func Info(msg string, fields ...logger.Field) {
	Get().Info(msg, fields...)
}

// Warn level information
func Warn(msg string, fields ...logger.Field) {
	Get().Warn(msg, fields...)
}

// Error level information
func Error(msg string, fields ...logger.Field) {
	Get().Error(msg, fields...)
}

// Panic level information
func Panic(msg string, fields ...logger.Field) {
	Get().Panic(msg, fields...)
}

// Fatal level information
func Fatal(msg string, fields ...logger.Field) {
	Get().Fatal(msg, fields...)
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
	err := Get().Sync()
	if err != nil && !strings.Contains(err.Error(), "/dev/stdout") {
		return err
	}
	return nil
}

// WithFields carrying field information
func WithFields(fields ...logger.Field) logger.ILoggerMethod {
	//return GetWithSkip(0).With(fields...)
	return Get().WithFields(fields...)
}
