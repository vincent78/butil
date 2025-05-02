package netUtil

import (
	"gitee.com/vincent78/gcutil/logger"
	"net"
)

// 常量
const (
	UDP = "udp"
)

// udp连接信息结构体
type UDPTransfer struct {
	// 连接对象 UDP这边的Conn是struct TCP那边的Conn是Interface
	Conn *net.UDPConn
	// 传输时的缓冲
	Buf [Buf_Len_1024 * 2]byte
}

func NewUDPTransfer(conn *net.UDPConn) *UDPTransfer {
	return &UDPTransfer{
		Conn: conn,
		Buf:  [2048]byte{},
	}
}

// 校验UDP地址格式是否正确
func ValidateUDPAddr(udpAddr string) (addr *net.UDPAddr, err error) {
	// 校验地址格式是否正确
	if addr, err = net.ResolveUDPAddr(UDP, udpAddr); err != nil {
		logger.Error("udp resolve ip `%s` failed, errMsg: %s", udpAddr, err.Error())
		return nil, err
	}
	return addr, nil
}

// 获取udpAddr对象
// addr: ip:port  eg:
//func GetUDPAddr(addr string) (udpAddr *net.UDPAddr, err error) {
//	if udpAddr, err = net.ResolveUDPAddr(UDP, addr); err != nil {
//		logger.Error("udp resolve ip `%s` failed, errMsg: %s", addr, err.Error())
//	}
//	return udpAddr, nil
//}
