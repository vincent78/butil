package netUtil

import (
	"bytes"
	"net"
	"testing"
)

func TestStartUDPServer(t *testing.T) {
	server, err := NewUDPServer("localhost:9009", nil, func(transfer *UDPTransfer) {
		flag := true
		defer transfer.ServerCloseConn()
		for flag {
			var (
				n          int
				err        error
				remoteAddr *net.UDPAddr
			)
			// 接收
			if n, remoteAddr, err = transfer.Conn.ReadFromUDP(transfer.Buf[:]); err != nil {
				t.Log(err)
				return
			}
			msg := transfer.Buf[:n]
			t.Logf("receive `%s` from client(%s)", msg, remoteAddr.String())
			// 发送
			reply := append([]byte("> "), msg...)
			if string(bytes.Trim(msg, "\r\n\t ")) == "quit" {
				reply = []byte("Bye!")
				flag = false
			}
			if n, err = transfer.Conn.WriteToUDP(reply, remoteAddr); err != nil {
				t.Log(err)
				return
			}
			t.Logf("send `%s` to client(%s)", reply, remoteAddr.String())
		}
	})
	if err != nil {
		t.Log(err)
	}
	server.Run()
}
