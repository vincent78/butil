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

	"github.com/vincent78/butil/config"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	formatConsole = "console"
	formatJSON    = "json"

	levelDebug = "DEBUG"
	levelInfo  = "INFO"
	levelWarn  = "WARN"
	levelError = "ERROR"
)

type Logger = zap.Logger
type SugaredLogger = zap.SugaredLogger

var defaultLogger *Logger
var defaultSugaredLogger *SugaredLogger
var loggerMap = make(map[string]*Logger)

func getLogger() *Logger {
	checkNil()
	return defaultLogger.WithOptions(zap.AddCallerSkip(1))
}

func getSugaredLogger() *SugaredLogger {
	checkNil()
	return defaultSugaredLogger.WithOptions(zap.AddCallerSkip(1))
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
func Init(opts ...Option) (*Logger, error) {
	o := defaultOptions()
	o.apply(opts...)
	isSave := o.isSave
	levelName := o.level
	encoding := o.encoding
	disableCaller := o.disableCaller

	var err error
	var zapLog *Logger
	var str string
	if !isSave {
		zapLog, err = log2Terminal(levelName, encoding, disableCaller)
		if err != nil {
			panic(err)
		}
		str = fmt.Sprintf("initialize logger finish, config is output to 'terminal', format=%s, level=%s", encoding, levelName)
	} else {
		zapLog = log2File(encoding, levelName, o.fileConfig)
		str = fmt.Sprintf("initialize logger finish, config is output to 'file', format=%s, level=%s, file=%s", encoding, levelName, o.fileConfig.filename)
	}

	if len(o.hooks) > 0 {
		zapLog = zapLog.WithOptions(zap.Hooks(o.hooks...))
	}
	if defaultLogger == nil {
		defaultLogger = zapLog
		defaultSugaredLogger = zapLog.Sugar()
	}
	zapLog.Info(str)
	return zapLog, err
}

func log2Terminal(levelName string, encoding string, disableCaller bool) (*Logger, error) {
	js := fmt.Sprintf(`{
      		"level": "%s",
            "encoding": "%s",
      		"outputPaths": ["stdout"],
            "errorOutputPaths": ["stdout"],
			"disableCaller": %v
		}`, levelName, encoding, disableCaller)

	var config zap.Config
	err := json.Unmarshal([]byte(js), &config)
	if err != nil {
		return nil, err
	}

	config.EncoderConfig = zap.NewProductionEncoderConfig()
	if encoding == formatConsole {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // logging color
	} else {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // logging levels in the log file using upper case letters
	}
	config.EncoderConfig.EncodeTime = timeFormatter // default time format
	return config.Build()
}

func log2File(encoding string, levelName string, fo *fileOptions) *Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // modify Time Encoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // logging levels in the log file using upper case letters
	var encoder zapcore.Encoder
	if encoding == formatConsole { // console format
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else { // json format
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fo.filename,      // file name
		MaxSize:    fo.maxSize,       // maximum file size (MB)
		MaxBackups: fo.maxBackups,    // maximum number of old files
		MaxAge:     fo.maxAge,        // maximum number of days for old documents
		Compress:   fo.isCompression, // whether to compress and archive old files
	})
	core := zapcore.NewCore(encoder, ws, getLevelSize(levelName))

	// add the function call information log to the log.
	return zap.New(core, zap.AddCaller())
}

// DEBUG(default), INFO, WARN, ERROR
func getLevelSize(levelName string) zapcore.Level {
	levelName = strings.ToUpper(levelName)
	switch levelName {
	case levelDebug:
		return zapcore.DebugLevel
	case levelInfo:
		return zapcore.InfoLevel
	case levelWarn:
		return zapcore.WarnLevel
	case levelError:
		return zapcore.ErrorLevel
	}
	return zapcore.DebugLevel
}

func timeFormatter(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// GetWithSkip get defaultLogger, set the skipped caller value, customize the number of lines of code displayed
func GetWithSkip(skip int) *Logger {
	checkNil()
	return defaultLogger.WithOptions(zap.AddCallerSkip(skip))
}

// Get logger
func Get() *Logger {
	checkNil()
	return defaultLogger
}

func checkNil() {
	if defaultLogger == nil {
		_, err := Init() // default output to console
		if err != nil {
			panic(err)
		}
	}
}

func InitLoggerByConfig(cfgs map[string]config.LoggerConfig) {
	for name, cfg := range cfgs {

		log, err := Init(
			WithLevel(cfg.Level),
			WithFormat(cfg.Format),
			WithSave(
				cfg.IsSave,
				WithFileName(cfg.LogFileConfig.Filename),
				WithFileMaxSize(cfg.LogFileConfig.MaxSize),
				WithFileMaxBackups(cfg.LogFileConfig.MaxBackups),
				WithFileMaxAge(cfg.LogFileConfig.MaxAge),
				WithFileIsCompression(cfg.LogFileConfig.IsCompression),
			),
		)
		if err != nil {
			panic("init logger error:" + err.Error())
		}

		loggerMap[name] = log
	}
}

func GetLogger(name string) *Logger {
	if logger, ok := loggerMap[name]; ok {
		return logger
	} else {
		panic("logger not found: " + name)
	}
}

func SetDefaultLogger(name string) {
	defaultLogger = GetLogger(name)
}

func InitByConf(cfg config.LoggerConfig) (*Logger, error) {
	// initializing log
	return Init(
		WithLevel(cfg.Level),
		WithFormat(cfg.Format),
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
