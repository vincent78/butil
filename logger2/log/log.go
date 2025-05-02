// Package log provides debug logging
package log

import (
	"github.com/vincent78/butil/logger2/logger"
	log "github.com/vincent78/butil/logger2/logger"
	"github.com/vincent78/butil/logger2/writer"
	"github.com/vincent78/butil/logger2/zap"
	"github.com/vincent78/butil/utils/fileUtil"
	"io"
	"os"
)

// SetupLogger 日志 cap 单位为kb
func SetupLogger(opts ...Option) logger.Logger {
	op := setDefault()
	for _, o := range opts {
		o(&op)
	}
	if !fileUtil.Exist(op.path) {
		err := fileUtil.MkDir(op.path)
		if err != nil {
			logger.Fatalf("create dir error: %s", err.Error())
		}
	}
	var err error
	var output io.Writer
	switch op.stdout {
	case "file":
		output, err = writer.NewFileWriter(
			writer.WithPath(op.path),
			writer.WithCap(op.cap<<10),
		)
		if err != nil {
			logger.Fatal("logger setup error: %s", err.Error())
		}
	default:
		output = os.Stdout
	}
	var level logger.Level
	level, err = logger.GetLevel(op.level)
	if err != nil {
		logger.Fatalf("get logger level error, %s", err.Error())
	}

	switch op.driver {
	case "zap":
		log.DefaultLogger, err = zap.NewLogger(logger.WithLevel(level), logger.WithOutput(output), zap.WithCallerSkip(2))
		if err != nil {
			logger.Fatalf("new zap logger error, %s", err.Error())
		}
	//case "logrus":
	//	setLogger = logrus.NewLogger(logger.WithLevel(level), logger.WithOutput(output), logrus.ReportCaller())
	default:
		log.DefaultLogger = logger.NewLogger(logger.WithLevel(level), logger.WithOutput(output))
	}
	return log.DefaultLogger
}
