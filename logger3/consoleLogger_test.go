package logger3

import (
	"context"
	"gitee.com/vincent78/gcutil/config"
	"gitee.com/vincent78/gcutil/utils/timeUtil"
	"testing"
	"time"
)

func TestNewConsoleLogger(t *testing.T) {
	log := NewConsoleLogger(context.Background(), config.LoggerConfig{
		Home:   "",
		Prefix: "",
		Level:  2,
	})

	ch := make(chan struct{})
	go func() {
		for i := 0; i < 5; i++ {
			log.Debug("this is the test %v - %v", timeUtil.NowStr(), i)
			time.Sleep(time.Millisecond * 500)
		}
		for i := 0; i < 5; i++ {
			log.Info("this is the test %v - %v", timeUtil.NowStr(), i)
			time.Sleep(time.Millisecond * 500)
		}
		for i := 0; i < 5; i++ {
			log.Warn("this is the test %v - %v", timeUtil.NowStr(), i)
			time.Sleep(time.Millisecond * 500)
		}
		for i := 0; i < 5; i++ {
			log.Error("this is the test %v - %v", timeUtil.NowStr(), i)
			time.Sleep(time.Millisecond * 500)
		}
		ch <- struct{}{}
	}()
	<-ch
}
