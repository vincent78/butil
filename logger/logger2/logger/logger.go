package logger

var (
	// DefaultLogger logger
	DefaultLogger Logger
)

// Logger is a generic logging interface
type Logger interface {
	// Init initialises options
	Init(options ...Option) error
	// Options The Logger options
	Options() Options
	// Fields set fields to always be logged
	Fields(fields map[string]any) Logger
	// Log writes a log entry
	Log(level Level, v ...any)
	// Logf writes a formatted log entry
	Logf(level Level, format string, v ...any)
	// String returns the name of logger
	String() string
}

func Init(opts ...Option) error {
	return DefaultLogger.Init(opts...)
}

func Fields(fields map[string]any) Logger {
	return DefaultLogger.Fields(fields)
}

func Log(level Level, v ...any) {
	DefaultLogger.Log(level, v...)
}

func Logf(level Level, format string, v ...any) {
	DefaultLogger.Logf(level, format, v...)
}

func String() string {
	return DefaultLogger.String()
}
