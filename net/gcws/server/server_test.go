package server

import (
	"testing"

	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/logger/logger1"
	"github.com/vincent78/butil/net/gcgin"
	"github.com/vincent78/butil/sys"
)

func initServerLog() {
	path := "/tmp/ws"

	conf := logger1.NewLogConfig()
	conf.Path = path
	logger1.NewLogger(conf)

	conf.ToConsole = false
	conf.SimpleFile = true
	conf.ShowStack = true
	conf.Name = global.LogFileWSSName
	logger1.NewLogger(conf)

	conf.Name = global.LogFileHttpName
	logger1.NewLogger(conf)
}

/*
go test -v *.go -test.run TestWSGCServer
*/
func TestWSGCServer(t *testing.T) {
	initServerLog()
	e := gcgin.InitEngine(true)
	c := make(chan struct{})
	InitWSByGin(e, "/ws", func(s *WSServer) {

	}, NewWSServerConfig())
	go gcgin.StartServer(":8101", c, e)
	<-c
	sys.BlockBySignal()
}
