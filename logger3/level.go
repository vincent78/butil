package logger3

type LoggerLeverl int

const (
	LevelDebug LoggerLeverl = iota + 1
	LevelInfo
	LevelWarn
	LevelError
)

func (l LoggerLeverl) str() string {
	if l == LevelError {
		return "ERROR"
	} else if l == LevelWarn {
		return "WARN"
	} else if l == LevelInfo {
		return "INFO"
	} else if l == LevelDebug {
		return "DEBUG"
	} else {
		return "UNKNOW"
	}
}
