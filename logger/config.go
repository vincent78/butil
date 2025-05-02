package logger

import (
	"gitee.com/vincent78/gcutil/dup"
)

type LogConfig struct {
	Path       string `json:"path"`
	Name       string `json:"name"`
	Level      int    `json:"level"`
	Override   bool   `json:"override"`
	ToConsole  bool   `json:"toConsole"`
	Split2Day  bool   `json:"split2Day"`  // 是否需要拆分成按日期分的不同文件
	SimpleFile bool   `json:"simpleFile"` // true 所有日志都写入单个文件， false 日志按等级拆分成多个文件
	ShowStack  bool   `json:"showStack"`  // 当错误时是否显示stack
}

func NewLogConfig() *LogConfig {
	return &LogConfig{
		Path:       "",
		Name:       "",
		Level:      0,
		Override:   true,
		ToConsole:  true,
		Split2Day:  false,
		SimpleFile: true,
		ShowStack:  false,
	}
}

func NewLoggerConfByPath(path string) *LogConfig {

	conf := NewLogConfig()
	conf.Path = path
	return conf
}

func NewLoggerConfByPathAndName(path, name string) *LogConfig {
	conf := NewLogConfig()
	conf.Path = path
	conf.Name = name
	return conf
}

func (l *LogConfig) Clone() *LogConfig {
	r := NewLogConfig()
	dup.Copy(l, r)
	return r
}
