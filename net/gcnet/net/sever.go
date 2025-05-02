package gcnet

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"gitee.com/vincent78/gcutil/logger"
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
	logger.Info("==============================================")
	logger.Info("tcp listen at : %v", lisAddr)
	logger.Info("==============================================")
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if err == io.EOF {
					return
				}
				logger.Error("server accept error:%v", err)
				continue
			}

			if conf.Version == "9" {
				go readMessageOld(conn)
			} else {
				l := NewServerListener(conn.RemoteAddr().String())
				logger.Info("==== connected: %v", l.LogPrefix)
				c <- l
				go readMessage(conn, l, conf)
				go writeMessage(conn, l, conf)
			}
		}
	}()
	return nil
}

func writeMessage(conn net.Conn, l *ServerListener, conf *NetConfig) {
	logger.Info("prepare write for: %v", l.LogPrefix)
	defer func() {
		logger.Info("close write for: %v", l.LogPrefix)
	}()

	for {
		select {
		case bytes := <-l.SendCH:

			_, err := conn.Write(bytes)
			if err != nil {
				logger.Error("<<--[%v] error: %x", l.LogPrefix, bytes)
			} else {
				logger.Debug("<<--[%v]: %x", l.LogPrefix, bytes)
			}
		case <-l.CloseCH:
			return
		}
	}
}

func readMessage(conn net.Conn, l *ServerListener, conf *NetConfig) {
	logger.Info("prepare read for: %v", l.LogPrefix)
	defer func() {
		logger.Info("close read for: %v", l.LogPrefix)
	}()
	reader := bufio.NewReader(conn)
	for {
		recvBuf := make([]byte, conf.ReadBuffSize)
		n, err := reader.Read(recvBuf) // recv data
		if err != nil {
			if err != io.EOF {
				l.ErrorCH <- err
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					logger.Error("read timeout:", err)
				}
			} else {
				l.Close()
				return
			}
		}
		if n > 0 {
			buf := recvBuf[:n]
			logger.Debug("-->>[%v]: %v", l.LogPrefix, hex.EncodeToString(buf))
			l.ReceiveCH <- buf
		}
	}
}
