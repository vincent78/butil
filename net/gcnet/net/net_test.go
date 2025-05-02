package gcnet

import (
	"gitee.com/vincent78/gcutil/logger"
	"net"
	"testing"
)

func TestNetServer(t *testing.T) {
	RegistNetHandler("echo", func(data interface{}, conn net.Conn) {
		logger.Debug("echo the data:%v", data)
	})

	ch := make(chan *ServerListener)
	Server(NewNetConfig("0.0.0.0", 7003), ch)
}
