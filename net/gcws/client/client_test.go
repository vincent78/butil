package client

import (
	"fmt"
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/logger"
	"github.com/vincent78/butil/sys"
	"github.com/vincent78/butil/utils/timeUtil"
	"testing"
	"time"
)

func initClientLog() {
	path := "/tmp/wsc"
	conf := logger.NewLogConfig()
	conf.Path = path
	logger.NewLogger(conf)

	conf.Name = global.LogFileWSCName
	logger.NewLogger(conf)

	conf.Name = global.LogFileHttpName
	logger.NewLogger(conf)
}

func GetWSServer() string {
	//ip, _ := gcnet.ExternalIP()
	return fmt.Sprintf("127.0.0.1:%v", 8101)
}

/*
go test -v *.go -test.run TestWSClient
*/
func TestWSClient(t *testing.T) {
	initClientLog()
	wsClient := NewWSClient(defaultServerSchema, GetWSServer(), DefaultServerPath)
	err := wsClient.Conn()
	defer wsClient.DisConn()
	if err != nil {
		t.Errorf("connect the url error:%v", err)
		return
	}
	t.Logf("connected")

	sys.BlockBySignal()
	t.Logf("exit")
}

func TestWSClientWithPort(t *testing.T) {
	initClientLog()
	wsClient := NewWSClient(defaultServerSchema, GetWSServer(), DefaultServerPath)
	err := wsClient.Conn()
	wsClient.DisConn()
	// 重新再连一次
	time.AfterFunc(1*time.Second, func() {
		go func() {
			wsClient1 := NewWSClient(defaultServerSchema, GetWSServer(), DefaultServerPath)
			err = wsClient1.Conn()
			if err != nil {
				t.Errorf("connect the url error:%v", err)
				return
			}
			t.Logf("connected")
			wsClient.DisConn()
		}()
	})
	//defer wsClient.DisConn()

	if err != nil {
		t.Errorf("connect the url error:%v", err)
		return
	}
	t.Logf("connected")

	sys.BlockBySignal()
	t.Logf("exit")
}

func TestWSClientConnWithLong(t *testing.T) {
	initClientLog()
	wsClient := NewWSClient(defaultServerSchema, GetWSServer(), DefaultServerPath)
	err := wsClient.Conn()
	defer wsClient.DisConn()
	if err != nil {
		t.Errorf("connect the url error:%v", err)
		tick := time.NewTicker(5 * time.Second)
		ctx := sys.SignalContext()
		for {
			select {
			case <-tick.C:
				//t.Logf("connect again : %v", timeUtil.NowStr())
				//err = wsClient.Conn()
				//if err != nil {
				//	t.Logf("connect success.")
				//	tick.Stop()
				//}

			case <-ctx.Done():
				t.Logf("exit")
				return
			}
		}
	}
}

func TestClientPing(t *testing.T) {
	initClientLog()
	wsClient := NewWSClient(defaultServerSchema, GetWSServer(), DefaultServerPath)
	err := wsClient.Conn()
	defer wsClient.DisConn()
	if err != nil {
		t.Errorf("connect the url error:%v", err)
		return
	}

	global.Timewheel.AddCron(10*time.Second, func() {
		logger.Error("time-wheel: %v", timeUtil.NowStr())
		wsClient.ping()
	})

	sys.BlockBySignal()
	t.Logf("exit")
}
