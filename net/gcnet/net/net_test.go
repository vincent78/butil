package gcnet

import (
	"net"
	"testing"

	"github.com/vincent78/butil/logger/logger1"
)

func TestNetServer(t *testing.T) {
	RegistNetHandler("echo", func(data interface{}, conn net.Conn) {
		logger1.Debug("echo the data:%v", data)
	})

	ch := make(chan *ServerListener)
	err := Server(NewNetConfig("0.0.0.0", 7003), ch)
	if err != nil {
		return
	}
}
