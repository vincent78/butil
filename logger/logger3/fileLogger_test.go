package logger3

import (
	"context"
	"github.com/vincent78/butil/config"
	"github.com/vincent78/butil/utils/timeUtil"
	"testing"
	"time"
)

func TestNewFileLogger(t *testing.T) {
	log, _ := NewFileLogger(context.Background(), config.LoggerConfig{
		Home:   "/tmp/test",
		Prefix: "",
		Level:  0,
	})

	ch := make(chan struct{})
	go func() {
		for i := 0; i < 10; i++ {
			log.Error("this is the test %v - %v", timeUtil.NowStr(), i)
			time.Sleep(500 * time.Millisecond)
		}
		ch <- struct{}{}
	}()
	<-ch
}
