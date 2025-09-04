package log

import (
	"testing"
	"time"

	"github.com/vincent78/butil/logger/logger2/logger"
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
