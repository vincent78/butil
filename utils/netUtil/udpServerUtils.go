package netUtil

import (
	"github.com/vincent78/butil/logger/logger1"
	"net"
)

type UDPServer struct {
	// UDP没有listener
	Addr        *net.UDPAddr // ip:port
	InitHandler func(*UDPTransfer) error
	Handler     func(*UDPTransfer)
}

// NewUDPServer 创建UDP服务并返回对象
func NewUDPServer(addr string, initHandler func(*UDPTransfer) error, handler func(transfer *UDPTransfer)) (*UDPServer, error) {
	logger1.Info("New UDPServer addr: %s", addr)
	// 校验地址格式是否正确
	if udpAddr, err := ValidateUDPAddr(addr); err != nil {
		return nil, err
	} else {
		// 返回UDPServer对象
		return &UDPServer{
			Addr:        udpAddr,
			InitHandler: initHandler,
			Handler:     handler,
		}, nil
	}
}

// 启动UDP服务
func (uServ *UDPServer) Run() (err error) {
	addrStr := uServ.Addr.String()
	logger1.Info("UDPServer Run addr: %s", addrStr)

	for {
		err = uServ.acceptAndHandle()
		if err != nil {
			break
		}
	}
	return
}

// 服务端断开conn
func (uTrans *UDPTransfer) ServerCloseConn() error {
	// 注 可以调用uTrans.Conn.RemoteAddr() 但返回值是nil
	if uTrans.Conn == nil {
		return nil
	}
	logger1.Error("UDPServer close Conn, ServerAddr: %s", uTrans.Conn.LocalAddr().String())
	if err := uTrans.Conn.Close(); err != nil {
		logger1.Error("UDPServer `%s` close conn failed, errMsg: %s", uTrans.Conn.LocalAddr().String(), err.Error())
		return err
	}
	return nil
}

// udp直接就是conn 没有所谓的listener 关闭服务器和上面的ServerCloseConn功能完全重叠
//func (uServ *UDPServer) CloseServer() (err error) {}

func (uServ *UDPServer) acceptAndHandle() error {
	defer func() {
		if tmp := recover(); tmp != nil {
			logger1.Error("udp panic: %v", tmp)
		}
		logger1.Info("====================== UDP SERVER RUN END ======================")
	}()

	logger1.Info("====================== UDP SERVER RUN START ======================")
	var (
		conn *net.UDPConn
		err  error
	)
	conn, err = net.ListenUDP(UDP, uServ.Addr)
	if err != nil {
		logger1.Error("udp: create conn failed, errMsg: %s", err.Error())
		return err
	}
	transfer := NewUDPTransfer(conn)
	// 初始化
	if uServ.InitHandler != nil {
		logger1.Info("udp: server( %s ) exec InitHandler...", uServ.Addr.String())
		if err = uServ.InitHandler(transfer); err != nil {
			transfer.ServerCloseConn()
			return nil
		}
	}
	// 执行正常流程
	if uServ.Handler != nil {
		logger1.Info("udp: server( %s ) exec Handler...", uServ.Addr.String())
		uServ.Handler(transfer)
	}

	return nil
}
