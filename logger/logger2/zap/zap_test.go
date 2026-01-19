package zap

import (
	"fmt"
	"testing"

	logger2 "github.com/vincent78/butil/logger/logger2/logger"
	"github.com/vincent78/butil/logger/logger2/writer"
)

func TestName(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}

	if l.String() != "zap" {
		t.Errorf("name is error %s", l.String())
	}

	t.Logf("test logger name: %s", l.String())
}

func TestLogf(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}

	logger2.DefaultLogger = l
	logger2.Logf(logger2.InfoLevel, "test logf: %s", "name")
}

func TestSetLevel(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}
	logger2.DefaultLogger = l

	logger2.Init(logger2.WithLevel(logger2.DebugLevel))
	l.Logf(logger2.DebugLevel, "test show debug: %s", "debug msg")

	logger2.Init(logger2.WithLevel(logger2.InfoLevel))
	l.Logf(logger2.DebugLevel, "test non-show debug: %s", "debug msg")
}

func TestWithReportCaller(t *testing.T) {
	var err error
	logger2.DefaultLogger, err = NewLogger(WithCallerSkip(0))
	if err != nil {
		t.Fatal(err)
	}

	logger2.Logf(logger2.InfoLevel, "testing: %s", "WithReportCaller")
}

func TestFields(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}
	logger2.DefaultLogger = l.Fields(map[string]any{
		"x-request-id": "123456abc",
	})
	logger2.DefaultLogger.Log(logger2.InfoLevel, "hello")
}

func TestFile(t *testing.T) {
	//output, err := writer.NewFileWriter("testdata", "log")
	output, err := writer.NewFileWriter()
	if err != nil {
		t.Errorf("logger setup error: %s", err.Error())
	}
	//var err error
	logger2.DefaultLogger, err = NewLogger(logger2.WithLevel(logger2.TraceLevel), WithOutput(output))
	if err != nil {
		t.Errorf("logger setup error: %s", err.Error())
	}
	logger2.DefaultLogger = logger2.DefaultLogger.Fields(map[string]any{
		"x-request-id": "123456abc",
	})
	fmt.Println(logger2.DefaultLogger)
	logger2.DefaultLogger.Log(logger2.InfoLevel, "hello")
}
