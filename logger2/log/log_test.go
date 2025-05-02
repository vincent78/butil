package log

import (
	"gitee.com/vincent78/gcutil/logger2/logger"
	"testing"
	"time"
)

func TestSetupLogger(t *testing.T) {
	log := SetupLogger(
		WithType(""),
		WithPath("/tmp/logTest"),
		WithLevel("trace"),
		WithStdout("file"),
		WithCap(0),
	)
	log.Logf(logger.TraceLevel, "test %v", time.Now())
}
