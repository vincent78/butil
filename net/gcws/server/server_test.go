package server

import (
	"gitee.com/vincent78/gcutil/global"
	"gitee.com/vincent78/gcutil/logger"
	"gitee.com/vincent78/gcutil/net/gcgin"
	"gitee.com/vincent78/gcutil/sys"
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
