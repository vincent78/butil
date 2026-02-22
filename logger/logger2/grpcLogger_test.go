package logger2

import (
	"testing"

	"github.com/vincent78/butil/logger/logger4"
	"go.uber.org/zap"
)

func TestReplaceGRPCLoggerV2(t *testing.T) {
	ReplaceGRPCLoggerV2(logger4.Get())

	l := &grpcLogger{
		zLog:      logger4.Get().WithOptions(zap.AddCallerSkip(1)),
		verbosity: 0,
	}

	l.V(0)

	l.Info("test info")
	l.Infof("test %s", "info")
	l.Infoln("test info")

	l.Warning("test warning")
	l.Warningf("test %s", "warning")
	l.Warningln("test warning")

	l.Error("test error")
	l.Errorf("test %s", "error")
	l.Errorln("test error")

	//l.Fatal("test fatal")
	//l.Fatalf("test %s", "fatal")
	//l.Fatalln("test fatal")
}
