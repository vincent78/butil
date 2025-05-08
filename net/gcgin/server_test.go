package gcgin

import (
	"fmt"
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/logger/logger1"
	"regexp"
	"testing"
)

func initLogger() {
	path := "/tmp/http"
	conf := logger1.NewLogConfig()
	conf.Path = path
	logger1.NewLogger(conf)

	conf = logger1.NewLogConfig()
	conf.Path = path
	conf.Name = global.LogFileHttpName
	conf.ToConsole = false
	logger1.NewLogger(conf)
}

func TestServer(t *testing.T) {
	initLogger()
	g := InitEngine(true)
	StartServer(":8101", nil, g)
}

func TestServerIntercept(t *testing.T) {
	initLogger()

	if ok, _ := regexp.Match("^/swagger", []byte("/swagger/test")); ok {
		fmt.Println("Match Found!")
	} else {
		fmt.Println("Not Match!")
	}

	if ok, _ := regexp.Match("^/swagger", []byte("//swagger/test")); ok {
		fmt.Println("Match Found!")
	} else {
		fmt.Println("Not Match!")
	}

	if ok, _ := regexp.Match("^/swagger", []byte("swagger/test")); ok {
		fmt.Println("Match Found!")
	} else {
		fmt.Println("Not Match!")
	}

	if ok, _ := regexp.Match("^/swagger", []byte("/swagger/")); ok {
		fmt.Println("Match Found!")
	} else {
		fmt.Println("Not Match!")
	}
}
