package logger1

import (
	"fmt"
)

func genNameByLevel(conf *LogConfig, level int) string {
	levelName := DefaultKeyError
	if !conf.SimpleFile {
		if level == LevelDebug {
			levelName = DefaultKeyDebug
		} else if level == LevelInfo {
			levelName = DefaultKeyInfo
		} else if level == LevelWarn {
			levelName = DefaultKeyWarn
		}
	}
	if conf.Name == "" {
		return DefaultKeyDebug
	} else {
		if conf.SimpleFile {
			return conf.Name
		} else {
			return fmt.Sprintf("%v-%v", levelName, levelName)
		}
	}
}

func Debug(msg string, args ...any) {
	writeByLevel(DefaultMapKey, LevelDebug, msg, args...)
}

func Info(msg string, args ...any) {
	writeByLevel(DefaultMapKey, LevelInfo, msg, args...)
}

func Warn(msg string, args ...any) {
	writeByLevel(DefaultMapKey, LevelWarn, msg, args...)
}

func Error(msg string, args ...any) {
	writeByLevel(DefaultMapKey, LevelError, msg, args...)
	//writeByLevel(DefaultMapKey, LevelError, "%v", strUtil.Bytes2String(debug.Stack()))
}

func DebugByName(name, msg string, args ...any) {
	if name != "" {
		writeByLevel(name, LevelDebug, msg, args...)
	} else {
		Debug(msg, args...)
	}
}

func InfoByName(name, msg string, args ...any) {
	if name != "" {
		writeByLevel(name, LevelInfo, msg, args...)
	} else {
		Info(msg, args...)
	}
}

func WarnByName(name, msg string, args ...any) {

	if name != "" {
		writeByLevel(name, LevelWarn, msg, args...)
	} else {
		Warn(msg, args...)
	}
}

func ErrorByName(name, msg string, args ...any) {

	if name != "" {
		writeByLevel(name, LevelError, msg, args...)
		//writeByLevel(name, LevelError, "%v", strUtil.Bytes2String(debug.Stack()))
	} else {
		Error(msg, args...)
	}
}

func writeByLevel(name string, level int, msg string, args ...any) {
	if name == "" {
		name = DefaultMapKey
	}
	defaultObj := loggerContainer[name]
	if defaultObj != nil && defaultObj.Channel != nil {
		fmtMsg := msg
		msgObj := LogMessage{
			Msg:   fmt.Sprintf(fmtMsg, args...),
			Level: level,
		}
		defaultObj.Channel <- msgObj
	} else {
		fmt.Println(fmt.Sprintf(msg, args...))
	}
}
