package logger

const DefaultPath = "./logs"

const DefaultMapKey = "default"

const DefaultChannelCache = 10

/**
level:
	zap.DebugLevel		-1
	zap.InfoLevel 		0
	zap.WarnLevel		1
	zap.ErrorLevel		2
	zap.DPanicLevel		3
	zap.PanicLevel		4
	zap.FatalLevel		5
*/

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

const (
	DefaultKeyDebug = "debug"
	DefaultKeyInfo  = "info"
	DefaultKeyWarn  = "warn"
	DefaultKeyError = "error"
)

const (
	MillTimeFormat = "2006-01-02 15:04:05.000"
)

var loggerContainer = map[string]*Logger{}

//var LogHome = DefaultPath
