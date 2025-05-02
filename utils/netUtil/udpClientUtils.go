package netUtil

import (
	"github.com/vincent78/butil/logger"
	"net"
)

type UDPClient struct {
	Addr    *net.UDPAddr       // ip:port
	Handler func(*UDPTransfer) // 每次收到服务端信息后，执行的处理信息函数
}

// 创建UDP 客户端连接对象并返回
// handler 函数，主要用于客户端收到服务器消息后的处理函数
// 返回 UDPClient 对象
func NewUDPClient(addr string, handler func(*UDPTransfer)) (*UDPClient, error) {
	// 校验地址格式是否正确
	if udpAddr, err := ValidateUDPAddr(addr); err != nil {
		return nil, err
	} else {
		return &UDPClient{
			Addr:    udpAddr,
			Handler: handler,
		}, nil
	}
}

// 客户端与服务器建连接
// 连接创建之后，进行回调处理服务器消息函数
func (uCli *UDPClient) Conn() error {
	if conn, err := net.DialUDP(UDP, nil, uCli.Addr); err != nil {
		logger.Error("connect server failed, errMsg: %s", err.Error())
		return err
	} else {
		transfer := NewUDPTransfer(conn)
		if uCli.Handler != nil {
			uCli.Handler(transfer)
		}
	}
	return nil
}

// 客户端关闭conn
func (uTrans *UDPTransfer) ClientCloseConn() error {
	if uTrans.Conn == nil {
		return nil
	}
	logger.Error("UDPClient close Conn, ClientAddr: %s", uTrans.Conn.LocalAddr().String())
	if err := uTrans.Conn.Close(); err != nil {
		logger.Error("UDPClient `%s` close conn failed, errMsg: %s", uTrans.Conn.LocalAddr().String(), err.Error())
	}
	return nil
}
