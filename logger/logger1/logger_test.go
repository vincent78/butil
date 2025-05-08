package logger1

import (
	"github.com/vincent78/butil/utils/strUtil"
	"sync"
	"testing"
	"time"
)

func TestCopyLoggerConfig(t *testing.T) {
	conf1 := NewLogConfig()
	t.Logf("conf no.1: %v", conf1)
	conf2 := conf1.Clone()
	conf2.Path = "testPath1"
	conf2.Name = "testName1"
	conf2.SimpleFile = !conf1.SimpleFile
	conf2.ToConsole = !conf1.ToConsole

	t.Logf("src: %v", strUtil.ToJsonStr(conf1))
	t.Logf("now: %v", strUtil.ToJsonStr(conf2))
}

func TestNewLogger(t *testing.T) {
	NewLogger(nil)
	Error("default %d", 1)
	time.Sleep(time.Second * 2)
}

func TestNewLogger2(t *testing.T) {
	NewLogger(nil)
	conf := NewLogConfig()
	conf.Name = "test"
	NewLogger(conf)
	ErrorByName("test", "error default %d", 1)
	time.Sleep(time.Second * 2)
}

func TestNewLogger3(t *testing.T) {
	NewLogger(nil)
	conf := NewLogConfig()
	conf.Name = "test"
	NewLogger(conf)
	for i := 0; i < 100; i++ {
		Error("--- %d", i)
		ErrorByName("test", "debug %d", i)
	}
	time.Sleep(time.Second * 2)
}

func TestNewLogger4Split2Day(t *testing.T) {
	conf := NewLogConfig()
	conf.Name = "test"
	conf.Split2Day = true
	NewLogger(conf)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		num := 0
		for {
			t := time.NewTimer(time.Second * 2)
			select {
			case tn := <-t.C:
				DebugByName("test", "%v", tn)
				num += 1
				if num < 200 {
					t.Reset(time.Second * 2)
				} else {
					return
				}
			}
		}
	}(wg)
	wg.Wait()
	t.Logf("---- end")
}
