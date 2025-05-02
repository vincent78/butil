package logger3

import (
	"context"
	"fmt"
	"gitee.com/vincent78/gcutil/config"
)

type LogModel struct {
	//Path  string       `json:"path"`
	//Name  string       `json:"name"`
	Level LoggerLeverl `json:"level"`
	//Override   bool         `json:"override"`
	//SimpleFile bool `json:"simpleFile"` // true 所有日志都写入单个文件， false 日志按等级拆分成多个文件
	//ShowStack  bool          `json:"showStack"`  // 当错误时是否显示stack
	MsgFormat MessageFormat `json:"msgFormat"`
	Queue     chan *MsgModel
	ctx       context.Context
}

type LogOperation interface {
	InitWithConfig(config config.LoggerConfig) error

	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})

	DebugWithField(msg string, field map[string]interface{})
	InfoWithField(msg string, field map[string]interface{})
	WarnWithField(msg string, field map[string]interface{})
	ErrorWithField(msg string, field map[string]interface{})
}

func (log *LogModel) Debug(msg string, args ...interface{}) {
	if log.Level <= LevelDebug {
		m := NewMsgModel(LevelDebug, fmt.Sprintf(msg, args...), nil)
		log.Queue <- m
	}
}
func (log *LogModel) Info(msg string, args ...interface{}) {
	if log.Level <= LevelInfo {
		m := NewMsgModel(LevelInfo, fmt.Sprintf(msg, args...), nil)
		log.Queue <- m
	}
}
func (log *LogModel) Warn(msg string, args ...interface{}) {
	if log.Level <= LevelWarn {
		m := NewMsgModel(LevelWarn, fmt.Sprintf(msg, args...), nil)
		log.Queue <- m
	}
}
func (log *LogModel) Error(msg string, args ...interface{}) {
	if log.Level <= LevelError {
		m := NewMsgModel(LevelError, fmt.Sprintf(msg, args...), nil)
		log.Queue <- m
	}
}

func (log *LogModel) DebugWithField(msg string, field map[string]interface{}) {
	if log.Level <= LevelDebug {
		m := NewMsgModel(LevelDebug, msg, field)
		log.Queue <- m
	}
}
func (log *LogModel) InfoWithField(msg string, field map[string]interface{}) {
	if log.Level <= LevelInfo {
		m := NewMsgModel(LevelInfo, msg, field)
		log.Queue <- m
	}
}
func (log *LogModel) WarnWithField(msg string, field map[string]interface{}) {
	if log.Level <= LevelWarn {
		m := NewMsgModel(LevelWarn, msg, field)
		log.Queue <- m
	}
}
func (log *LogModel) ErrorWithField(msg string, field map[string]interface{}) {
	if log.Level <= LevelError {
		m := NewMsgModel(LevelError, msg, field)
		log.Queue <- m
	}
}
