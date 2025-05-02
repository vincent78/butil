package gchttp

import (
	"gitee.com/vincent78/gcutil/logger"
	"gitee.com/vincent78/gcutil/model"
	"testing"
)

func initLogger() {
	path := "/tmp"
	file := "test"

	c1 := logger.NewLogConfig()
	c1.Path = path
	c1.Name = file
	logger.NewLogger(c1)
}

func TestHttp(t *testing.T) {
	initLogger()
	c := make(chan struct{})
	go StartServer(":9090", c)
	<-c
	r := GetRequest("http://localhost:9090/timestamp")
	if r.Code == model.Success {
		logger.Debug("http get result. %+v ", r.Data)
	} else {
		logger.Error("http get error: %+v", r.Message)
	}
}
