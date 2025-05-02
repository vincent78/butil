package netUtil

import (
	"gitee.com/vincent78/gcutil/logger"
	"net"
)

type TCPClient struct {
	Addr    *net.TCPAddr       // ip:port
	Handler func(*TCPTransfer) // 每次收到服务端信息后，执行的处理信息函数
}

// 创建TCP 客户端连接对象并返回
// handler 函数，主要用于客户端收到服务器消息后的处理函数
// 返回 TCPClient 对象
func NewTCPClient(addr string, handler func(transfer *TCPTransfer)) (t *TCPClient, err error) {
	// 校验地址格式是否正确
	var tcpAddr *net.TCPAddr
	if tcpAddr, err = ValidateTCPAddr(addr); err != nil {
		return t, err
	}
	return &TCPClient{
		Addr:    tcpAddr,
		Handler: handler,
	}, nil
}

// 客户端与服务器建连接
// 连接创建之后，进行回调处理服务器消息函数
func (t *TCPClient) Conn() (transfer *TCPTransfer, err error) {
	var conn net.Conn
	// 与服务器建立连接
	if conn, err = net.DialTCP(TCP, nil, t.Addr); err != nil {
		logger.Error("connect server failed, err: %v", err)
		return
	}
	transfer = NewTCPTransfer(conn)
	if t.Handler != nil {
		t.Handler(transfer)
	}
	return
}

//与服务器断开连接
func (t *TCPTransfer) CloseConn() (err error) {
	logger.Info("TCPTransfer close conn")
	if t == nil {
		return
	}
	if t.ConnState {
		if err := t.Conn.Close(); err != nil {
			logger.Error("close client conn: %v error,errMsg: %v", t.Conn.LocalAddr().String(), err)
		}
		t.ConnState = false
	}
	return
}
