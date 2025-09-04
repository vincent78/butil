// Package log provides debug logging
package log

import (
	"io"
	"os"

	logger2 "github.com/vincent78/butil/logger/logger2/logger"
	writer2 "github.com/vincent78/butil/logger/logger2/writer"
	zap2 "github.com/vincent78/butil/logger/logger2/zap"
	"github.com/vincent78/butil/utils/fileUtil"
)

// SetupLogger 日志 cap 单位为kb
func SetupLogger(opts ...Option) logger2.Logger {
	op := setDefault()
	for _, o := range opts {
		o(&op)
	}
	if !fileUtil.Exist(op.path) {
		err := fileUtil.MkDir(op.path)
		if err != nil {
			logger2.Fatalf("create dir error: %s", err.Error())
		}
	}
	var err error
	var output io.Writer
	switch op.stdout {
	case "file":
		output, err = writer2.NewFileWriter(
			writer2.WithPath(op.path),
			writer2.WithCap(op.cap<<10),
		)
		if err != nil {
			logger2.Fatal("logger setup error: %s", err.Error())
		}
	default:
		output = os.Stdout
	}
	var level logger2.Level
	level, err = logger2.GetLevel(op.level)
	if err != nil {
		logger2.Fatalf("get logger level error, %s", err.Error())
	}

	switch op.driver {
	case "zap":
		logger2.DefaultLogger, err = zap2.NewLogger(logger2.WithLevel(level), logger2.WithOutput(output), zap2.WithCallerSkip(2))
		if err != nil {
			logger2.Fatalf("new zap logger error, %s", err.Error())
		}
	//case "logrus":
	//	setLogger = logrus.NewLogger(logger.WithLevel(level), logger.WithOutput(output), logrus.ReportCaller())
	default:
		logger2.DefaultLogger = logger2.NewLogger(logger2.WithLevel(level), logger2.WithOutput(output))
	}
	return logger2.DefaultLogger
}
