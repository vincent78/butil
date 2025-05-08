package gcnet

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/vincent78/butil/logger/logger1"
	"io"
	"net"
)

func Server(conf *NetConfig, c chan *ServerListener) error {
	//监听
	lisAddr := fmt.Sprintf("%v:%v", conf.Host, conf.Port)
	listener, err := net.Listen(conf.Protocol, lisAddr)
	if err != nil {
		return err
	}
	logger1.Info("==============================================")
	logger1.Info("tcp listen at : %v", lisAddr)
	logger1.Info("==============================================")
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if err == io.EOF {
					return
				}
				logger1.Error("server accept error:%v", err)
				continue
			}

			if conf.Version == "9" {
				go readMessageOld(conn)
			} else {
				l := NewServerListener(conn.RemoteAddr().String())
				logger1.Info("==== connected: %v", l.LogPrefix)
				c <- l
				go readMessage(conn, l, conf)
				go writeMessage(conn, l, conf)
			}
		}
	}()
	return nil
}

func writeMessage(conn net.Conn, l *ServerListener, conf *NetConfig) {
	logger1.Info("prepare write for: %v", l.LogPrefix)
	defer func() {
		logger1.Info("close write for: %v", l.LogPrefix)
	}()

	for {
		select {
		case bytes := <-l.SendCH:

			_, err := conn.Write(bytes)
			if err != nil {
				logger1.Error("<<--[%v] error: %x", l.LogPrefix, bytes)
			} else {
				logger1.Debug("<<--[%v]: %x", l.LogPrefix, bytes)
			}
		case <-l.CloseCH:
			return
		}
	}
}

func readMessage(conn net.Conn, l *ServerListener, conf *NetConfig) {
	logger1.Info("prepare read for: %v", l.LogPrefix)
	defer func() {
		logger1.Info("close read for: %v", l.LogPrefix)
	}()
	reader := bufio.NewReader(conn)
	for {
		recvBuf := make([]byte, conf.ReadBuffSize)
		n, err := reader.Read(recvBuf) // recv data
		if err != nil {
			if err != io.EOF {
				l.ErrorCH <- err
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					logger1.Error("read timeout:", err)
				}
			} else {
				l.Close()
				return
			}
		}
		if n > 0 {
			buf := recvBuf[:n]
			logger1.Debug("-->>[%v]: %v", l.LogPrefix, hex.EncodeToString(buf))
			l.ReceiveCH <- buf
		}
	}
}
