package netUtil

import (
	"github.com/vincent78/butil/logger"
	"net"
	"time"
)

// 常量
const (
	TCP          = "tcp"
	Buf_Len_1024 = 1024
)

// tcp连接信息结构体
type TCPTransfer struct {
	// 连接对象
	Conn net.Conn
	// 传输时的缓冲
	Buf [Buf_Len_1024 * 2]byte
	// 连接关闭标识: true: 连接开启， false : 连接关闭
	ConnState bool
	// 连接唯一标识
	ConnKey string
	// 客户端上一次活动时间(单位：ms)
	ClientLastActiveTime int64
}

func NewTCPTransfer(conn net.Conn) *TCPTransfer {
	return &TCPTransfer{
		Conn:                 conn,
		Buf:                  [2048]byte{},
		ConnState:            true,
		ConnKey:              conn.RemoteAddr().String(),
		ClientLastActiveTime: time.Now().UnixMilli(),
	}
}

// 校验TCP地址格式是否正确
func ValidateTCPAddr(tcpAddr string) (addr *net.TCPAddr, err error) {
	// 校验地址格式是否正确
	if addr, err = net.ResolveTCPAddr(TCP, tcpAddr); err != nil {
		logger.Error("tcp resolve ip %v failed,errMsg: %v", tcpAddr, err)
		return nil, err
	}
	return addr, nil
}

// 获取tcpAddr 对象
// addr: ip:port  eg: 127.0.0.1:8000
func GetTCPAddr(addr string) (*net.TCPAddr, error) {
	var tcpAddr *net.TCPAddr
	var err error
	if tcpAddr, err = net.ResolveTCPAddr(TCP, addr); err != nil {
		logger.Error("tcp resolve ip : %v failed,err: %v", addr, err)
		return tcpAddr, err
	}
	return tcpAddr, err
}
