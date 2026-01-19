package logger

import (
	"fmt"
	"os"
)

type Level int8

const (
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel Level = iota - 2
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// InfoLevel is the default logging priority.
	// General operational entries about what's going on inside the application.
	InfoLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	ErrorLevel
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. highest level of severity.
	FatalLevel
)

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	}
	return ""
}

// LevelForGorm 转换成gorm日志级别
func (l Level) LevelForGorm() int {
	switch l {
	case FatalLevel, ErrorLevel:
		return 2
	case WarnLevel:
		return 3
	case InfoLevel, DebugLevel, TraceLevel:
		return 4
	default:
		return 1
	}
}

// Enabled returns true if the given level is at or above this level.
func (l Level) Enabled(lvl Level) bool {
	return lvl >= l
}

// GetLevel converts a level string into a logger Level value.
// returns an error if the input string does not match known values.
func GetLevel(levelStr string) (Level, error) {
	switch levelStr {
	case TraceLevel.String():
		return TraceLevel, nil
	case DebugLevel.String():
		return DebugLevel, nil
	case InfoLevel.String():
		return InfoLevel, nil
	case WarnLevel.String():
		return WarnLevel, nil
	case ErrorLevel.String():
		return ErrorLevel, nil
	case FatalLevel.String():
		return FatalLevel, nil
	}
	return InfoLevel, fmt.Errorf("Unknown Level String: '%s', defaulting to InfoLevel", levelStr)
}

func Info(args ...any) {
	DefaultLogger.Log(InfoLevel, args...)
}

func Infof(template string, args ...any) {
	DefaultLogger.Logf(InfoLevel, template, args...)
}

func Trace(args ...any) {
	DefaultLogger.Log(TraceLevel, args...)
}

func Tracef(template string, args ...any) {
	DefaultLogger.Logf(TraceLevel, template, args...)
}

func Debug(args ...any) {
	DefaultLogger.Log(DebugLevel, args...)
}

func Debugf(template string, args ...any) {
	DefaultLogger.Logf(DebugLevel, template, args...)
}

func Warn(args ...any) {
	DefaultLogger.Log(WarnLevel, args...)
}

func Warnf(template string, args ...any) {
	DefaultLogger.Logf(WarnLevel, template, args...)
}

func Error(args ...any) {
	DefaultLogger.Log(ErrorLevel, args...)
}

func Errorf(template string, args ...any) {
	DefaultLogger.Logf(ErrorLevel, template, args...)
}

func Fatal(args ...any) {
	DefaultLogger.Log(FatalLevel, args...)
	os.Exit(1)
}

func Fatalf(template string, args ...any) {
	DefaultLogger.Logf(FatalLevel, template, args...)
	os.Exit(1)
}

// Returns true if the given level is at or lower the current logger level
func V(lvl Level, log Logger) bool {
	l := DefaultLogger
	if log != nil {
		l = log
	}
	return l.Options().Level <= lvl
}
