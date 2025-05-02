package server

import (
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/logger"
	"github.com/vincent78/butil/net/gcgin"
	"github.com/vincent78/butil/sys"
	"testing"
)

func initServerLog() {
	path := "/tmp/ws"

	conf := logger.NewLogConfig()
	conf.Path = path
	logger.NewLogger(conf)

	conf.ToConsole = false
	conf.SimpleFile = true
	conf.ShowStack = true
	conf.Name = global.LogFileWSSName
	logger.NewLogger(conf)

	conf.Name = global.LogFileHttpName
	logger.NewLogger(conf)
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
