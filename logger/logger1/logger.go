package logger1

import (
	"context"
	"fmt"

	"github.com/vincent78/butil/utils/fileUtil"
	"go.uber.org/zap"
)

type Logger struct {
	Key        string             // the key of map
	Ctx        context.Context    // the context of log
	CancelFunc context.CancelFunc // the func of cancel logger
	zapLogger  *zap.Logger        // the real logger
	Channel    chan LogMessage    // the input chan
	Conf       *LogConfig
}

func (obj *Logger) prepare() error {
	path := obj.Conf.Path
	if has := fileUtil.Exist(path); !has {
		err := fileUtil.MkDir(path)
		if err != nil {
			return err
		}
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	obj.Ctx = ctx
	obj.CancelFunc = cancelFunc

	obj.zapLogger = genZapLogger(obj)

	obj.Channel = make(chan LogMessage, DefaultChannelCache)
	go func(obj *Logger) {
		for {
			//TODO:这里可以进行按日期进行拆分
			select {
			case msg := <-obj.Channel:
				if msg.Level == LevelError {
					obj.zapLogger.Error(msg.Msg)
				} else if msg.Level == LevelWarn {
					obj.zapLogger.Warn(msg.Msg)
				} else if msg.Level == LevelInfo {
					obj.zapLogger.Info(msg.Msg)
				} else {
					obj.zapLogger.Debug(msg.Msg)
				}
			case <-ctx.Done():
				return
			}
		}
	}(obj)
	return nil
}

func (obj *Logger) Destroy() {
	if obj.CancelFunc != nil {
		obj.CancelFunc()
	}
	obj.zapLogger = nil
}

func NewLogger(conf *LogConfig) *Logger {

	if conf == nil {
		conf = NewLogConfig()
	} else {
		// 这里拷贝一份配置，指针操作要注意
		conf = conf.Clone()
	}

	logPath := conf.Path
	if logPath == "" {
		logPath = DefaultPath
		conf.Path = DefaultPath
	}

	key := conf.Name
	if key == "" {
		key = DefaultMapKey
	}

	existObj := loggerContainer[key]
	if conf.Override && existObj != nil {
		existObj.Destroy()
	}

	obj := &Logger{
		Conf: conf,
	}

	err := obj.prepare()
	if err != nil {
		fmt.Printf("prepare the logger error: %v", err)
		return nil
	}

	loggerContainer[key] = obj

	return obj
}

func (obj *Logger) Debug(msg string, args ...any) {
	writeByLevel(obj.Conf.Name, LevelDebug, msg, args...)
}

func (obj *Logger) Info(msg string, args ...any) {
	writeByLevel(obj.Conf.Name, LevelInfo, msg, args...)
}

func (obj *Logger) Warn(msg string, args ...any) {
	writeByLevel(obj.Conf.Name, LevelWarn, msg, args...)
}

func (obj *Logger) Error(msg string, args ...any) {
	writeByLevel(obj.Conf.Name, LevelError, msg, args...)
}
