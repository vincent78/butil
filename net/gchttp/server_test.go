package gchttp

import (
	"context"
	"testing"

	logger "github.com/vincent78/butil/logger/logger1"
	"github.com/vincent78/butil/model"
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
	urlstr := "http://localhost:9090/timestamp"
	r := Done(NewTask(context.Background(), urlstr, nil), GetClient(urlstr))
	if r.Code == model.Success {
		logger.Debug("http get result. %v ", r.Data)
	} else {
		logger.Error("http get error: %v", r.Message)
	}
}
