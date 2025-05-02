package netUtil

import (
	"testing"
)

func TestStartTCPServer(t *testing.T) {
	initlogger()
	//logger.InitLogger("./log", "test.log", 0)
	server, err := NewTCPServer("localhost:9009", 1, -1, nil, func(transfer *TCPTransfer) {
		for {
			var len int
			var err error
			if len, err = transfer.Conn.Read(transfer.Buf[:]); err != nil {
				t.Log(err)
				return
			}
			msg := transfer.Buf[:len]
			t.Logf("receive client msg:%v", string(msg))
		}
	}, nil)
	if err != nil {
		t.Log(err)
	}
	server.Run()
	for {

	}
}
