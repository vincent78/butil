// Package logger is log library encapsulated in https://github.com/uber-go/zap
//
// Support for terminal printing and log saving.
// Support for automatic log file cutting.
// Support for json format and console log format output.
// Supports Debug, Info, Warn, Error, Panic, Fatal, also supports fmt.Printf-like log printing, Debugf, Infof, Warnf, Errorf, Panicf, Fatalf.
package logger4

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/vincent78/butil/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type SimpleLogger = zap.Logger
type SugaredLogger = zap.SugaredLogger

func getLoggerWithOptions(opts ...Option) *SimpleLogger {
	checkNil()
	l, err := Init(opts...)
	if err != nil {
		panic(err)
	}
	return l
}

func getSugaredLogger() *SugaredLogger {
	checkNil()
	//return defaultSugaredLogger.WithOptions(zap.AddCallerSkip(1))
	return defaultSugaredLogger
}

func getSugaredLoggerWithOptions(opts ...Option) *SugaredLogger {
	checkNil()
	l, err := Init(opts...)
	if err != nil {
		panic(err)
	}
	return l.Sugar()
}

// Init initial log settings
// print the debug level log in the terminal, example: Init()
// print the info level log in the terminal, example: Init(WithLevel("info"))
// print the json format, debug level log in the terminal, example: Init(WithFormat("json"))
// log with hooks, example: Init(WithHooks(func(zapcore.Entry) error{return nil}))
// output the log to the file out.log, using the default cut log-related parameters, debug-level log, example: Init(WithSave())
// output the log to the specified file, custom set the log file cut log parameters, json format, debug level log, example:
// Init(
//
//	  WithFormat("json"),
//	  WithSave(true,
//
//			WithFileName("my.log"),
//			WithFileMaxSize(5),
//			WithFileMaxBackups(5),
//			WithFileMaxAge(10),
//			WithFileIsCompression(true),
//		))
func Init(opts ...Option) (*SimpleLogger, error) {
	o := defaultOptions()
	o.apply(opts...)

	var err error
	var zapLog *SimpleLogger
	var str string
	if !o.isSave {
		zapLog, err = log2Terminal(o)
		if err != nil {
			panic(err)
		}
		str = fmt.Sprintf("initialize logger finish, config is output to 'terminal', format=%s, level=%s", o.encoding, LevelString(o.level))
	} else {
		zapLog = log2File(o)
		str = fmt.Sprintf("initialize logger finish, config is output to 'file', format=%s, level=%s, file=%s", o.encoding, LevelString(o.level), o.fileConfig.filename)
	}

	if len(o.hooks) > 0 {
		zapLog = zapLog.WithOptions(zap.Hooks(o.hooks...))
	}
	//if defaultLogger == nil {
	//	defaultLogger = zapLog
	//	defaultSugaredLogger = zapLog.Sugar()
	//}
	zapLog.Info(str)
	return zapLog, err
}

func log2Terminal(o *options) (*SimpleLogger, error) {
	js := fmt.Sprintf(`{
      		"level": "%s",
            "encoding": "%s",
      		"outputPaths": ["stdout"],
            "errorOutputPaths": ["stdout"],
			"x": %v
		}`, LevelString(o.level), o.encoding, o.disableCaller)

	var config zap.Config
	err := json.Unmarshal([]byte(js), &config)
	if err != nil {
		return nil, err
	}

	config.EncoderConfig = zap.NewProductionEncoderConfig()
	if o.encoding == formatConsole {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // logging color
	} else {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // logging levels in the log file using upper case letters
	}
	config.EncoderConfig.EncodeTime = timeFormatter // default time format
	return config.Build(zap.AddStacktrace(zapcore.Level(o.stacktrace)),
		zap.AddCallerSkip(o.callerSkip),
	)
}

func log2File(o *options) *SimpleLogger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // modify Time Encoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // logging levels in the log file using upper case letters
	var encoder zapcore.Encoder
	if o.encoding == formatConsole { // console format
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else { // json format
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   o.fileConfig.filename,      // file name
		MaxSize:    o.fileConfig.maxSize,       // maximum file size (MB)
		MaxBackups: o.fileConfig.maxBackups,    // maximum number of old files
		MaxAge:     o.fileConfig.maxAge,        // maximum number of days for old documents
		Compress:   o.fileConfig.isCompression, // whether to compress and archive old files
	})
	core := zapcore.NewCore(encoder, ws, zapcore.Level(o.level))

	// add the function call information log to the log.
	return zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.Level(o.stacktrace)),
		zap.AddCallerSkip(o.callerSkip),
	)
}

// DEBUG(default), INFO, WARN, ERROR
func getLevel(name string) zapcore.Level {
	level, err := zapcore.ParseLevel(strings.ToLower(name))
	if err != nil {
		panic(err)
	}
	return level
}

func LevelString(level Level) string {
	l := zapcore.Level(level)
	return l.String()
}

func timeFormatter(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func checkNil() {
	if defaultLogger == nil {
		_, err := Init() // default output to console
		if err != nil {
			panic(err)
		}
	}
}

func InitLoggerByConfs(cfgs map[string]config.LoggerConfig) {
	for name, cfg := range cfgs {
		configs[name] = cfg
		log, err := InitLoggerByConf(cfg)
		if err != nil {
			panic("init logger error:" + err.Error())
		}
		loggerMap[name] = &NormalLogger{
			logger: log,
			conf:   cfg,
		}
	}
}

func InitLoggerByConf(cfg config.LoggerConfig) (*SimpleLogger, error) {
	return Init(
		WithLevel(cfg.Level),
		WithFormat(cfg.Format),
		WithCaller(cfg.DisableCaller, cfg.CallerSkip),
		WithStacktraceLevel(cfg.StacktraceLevel),
		WithSave(
			cfg.IsSave,
			WithFileName(cfg.LogFileConfig.Filename),
			WithFileMaxSize(cfg.LogFileConfig.MaxSize),
			WithFileMaxBackups(cfg.LogFileConfig.MaxBackups),
			WithFileMaxAge(cfg.LogFileConfig.MaxAge),
			WithFileIsCompression(cfg.LogFileConfig.IsCompression),
		),
	)
}

// getCallerInfo 获取调用者的文件路径和行号（去掉根目录前缀）
func getCallerInfo() string {
	// 从调用栈中查找第一个非logger包的调用者
	for i := 1; i <= 10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}

		funcName := fn.Name()

		// 跳过logger包内部的函数
		if !strings.Contains(funcName, "logger/logger4") {
			// 动态获取项目根目录，去掉根目录前缀
			relativeFile := trimProjectRoot(file)
			// 合并文件路径和行号
			location := fmt.Sprintf("%s:%d", relativeFile, line)
			return location
		}
	}

	return "unknown:0"
}

func trimProjectRoot(filePath string) string {
	projectRoot := getProjectRoot()
	if projectRoot == "" {
		return filePath
	}
	if strings.HasPrefix(filePath, projectRoot) {
		relativePath := strings.TrimPrefix(filePath, projectRoot)
		return strings.TrimPrefix(relativePath, "/")
	}

	return filePath
}

// getProjectRoot 获取项目根目录
func getProjectRoot() string {
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return ""
}
