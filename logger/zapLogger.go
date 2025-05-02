package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"time"
)

func genZapLogger(obj *Logger) *zap.Logger {

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // 小写编码器
		EncodeTime:     normalTimeEncoder,           // yyyy-MM-dd HH:mi:ss.SSS
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
	}
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	cores := make([]zapcore.Core, 0, 4)

	i := 0
	max := 4

	if obj.Conf.SimpleFile {
		// 确认只有当前Level的文件存在
		i = obj.Conf.Level
		max = obj.Conf.Level + 1
	}

	for ; i < max; i++ {
		filePath := path.Join(obj.Conf.Path, fmt.Sprintf("%v.log", genNameByLevel(obj.Conf, i)))

		loggerWriter := genWriter(filePath)
		if obj.Conf.Split2Day {
			go func(loggerWriter *lumberjack.Logger) {
				for {
					nowTime := time.Now()
					nowTimeStr := nowTime.Format("2006-01-02")
					//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
					t2, _ := time.ParseInLocation("2006-01-02", nowTimeStr, time.Local)
					// 第二天零点时间戳
					next := t2.AddDate(0, 0, 1)
					after := next.UnixNano() - nowTime.UnixNano() - 1
					select {
					case <-time.After(time.Duration(after) * time.Nanosecond):
						err := loggerWriter.Rotate()
						if err != nil {
							fmt.Printf("\n\nlogger rotate error: %v   \n\n", err)
						}
					case <-obj.Ctx.Done():
						return
					}
				}
			}(loggerWriter)
		}

		cores = append(cores, zapcore.NewCore(consoleEncoder,
			zapcore.AddSync(loggerWriter),
			genEnablerFunc(i),
		))
	}

	if obj.Conf.ToConsole {
		core := zapcore.NewCore(consoleEncoder,
			zapcore.AddSync(os.Stdout),
			genEnablerFunc(obj.Conf.Level),
		)
		cores = append(cores, core)
	}

	core := zapcore.NewTee(cores...)
	//return zap.New(core)
	if obj.Conf.ShowStack {
		stack := zap.AddStacktrace(zap.ErrorLevel)
		return zap.New(core, stack)
		//caller := zap.AddCaller()
		//development := zap.Development()
		//return zap.New(core, stack, caller ,development)
	} else {
		return zap.New(core)
	}

}

func genWriter(filePath string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filePath, // 日志文件路径  "/tmp/logs/debug.log"
		MaxSize:    1024,     // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,       // 日志文件最多保存多少个备份
		MaxAge:     30,       // 文件最多保存多少天
		Compress:   true,     // 是否压缩
		LocalTime:  true,
	}
}

func genEnablerFunc(i int) zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		return lvl >= zapcore.Level(i-1)
	}
}

func normalTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}
	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, MillTimeFormat)
		return
	}
	enc.AppendString(t.Format(MillTimeFormat))
}
