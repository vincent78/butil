package netUtil

import (
	"errors"
	"fmt"
	"gitee.com/vincent78/gcutil/logger"
	"gitee.com/vincent78/gcutil/utils/concurrentUtil"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	connMapPrintPrefix = "cur server clientMap:[ "
	transferStrFmt     = "tcpTransfer%v:[connKey:%v, connState:%v, clientLastActiveTime:%v, hasExists:%v (ms)];"
	connMapPrintSuffix = "] "
)

type TCPServer struct {
	listener               net.Listener             // 服务监听
	Addr                   *net.TCPAddr             // ip:port
	InitHandler            func(*TCPTransfer) error // 第一次收到客户端信息后，执行的回执函数
	Handler                func(*TCPTransfer)       // 每次收到客户端信息后，执行的处理信息函数
	PanicHandler           func(interface{})        // 接收客户端请求 传输处理数据过程中遇到panic时的处理函数
	ClientConMap           sync.Map                 // 客户端传输对象: key 为客户端addr, value 为 *TCPTransfer
	ClientMaxNum           int                      // 允许连接的最大客户端数
	ServerState            bool                     // 服务状态： true ： 服务存活 ， false : 服务关闭
	clientConnTimeOut      int64                    // 服务器客户端连接超时时间,默认30s
	clientConnTimeOutState bool                     // 服务器监听客户端连接超时标识，默认不开启
}

const Con_Time_Out int64 = 30000

// NewTCPServer 创建TCP服务并返回对象
// param: addr tcp服务地址
// param: clientMaxNum 允许最大的客户端连接数
// param: clientConnTimeOut 客户连接超时时间 , -1 表示不限制， 0表示默认，其它数值表示超时时间，单位是毫秒
// param: initHandler 客户端连接建立后，初始化化函数
// param: handler  客户端逻辑处理函数
// param: panicHandler  panic时处理函数，传空则使用默认panic 处理函数
// return t tcpServer 对象
func NewTCPServer(addr string, clientMaxNum int, clientConnTimeOut int64, initHandler func(transfer *TCPTransfer) error, handler func(transfer *TCPTransfer), panicHandler func(interface{})) (t *TCPServer, err error) {
	logger.Info("New TCPServer addr:%v", addr)
	// 校验地址格式是否正确
	var tcpAddr *net.TCPAddr
	if tcpAddr, err = ValidateTCPAddr(addr); err != nil {
		return t, err
	}
	if clientConnTimeOut == 0 {
		clientConnTimeOut = Con_Time_Out
	}

	// 返回TCPServer 对象
	return &TCPServer{
		Addr:              tcpAddr,
		InitHandler:       initHandler,
		Handler:           handler,
		PanicHandler:      panicHandler,
		ClientMaxNum:      clientMaxNum,
		clientConnTimeOut: clientConnTimeOut,
	}, nil
}

// 启动TCP 服务
func (t *TCPServer) Run() (err error) {
	logger.Info("tcp: TCPServer Run addr: %s", t.Addr.String())
	t.listener, err = net.Listen(TCP, t.Addr.String())
	if err != nil {
		logger.Error("tcp: run tcp server( %s ) failed, errMsg: %s", t.Addr.String(), err.Error())
		return
	}

	go func(server *TCPServer) {
		defer func() {
			server.CloseServer()
			logger.Info("++++++++++++++++ PROXY TCP SERVICE END ++++++++++++++++")
		}()
		for {
			if server.acceptAndHandle() != nil {
				break
			}
		}
	}(t)

	return
}

// 关闭服务器监听
func (t *TCPServer) CloseServer() (err error) {
	if t == nil {
		return
	}
	logger.Info("TCPServer listener close start, addr:%v", t.Addr.String())
	if t.listener == nil {
		logger.Info("TCPServer listener not exit, close finish")
		return
	}
	// 先关闭连接
	logger.Info("TCPServer listener close: first close client conn")
	transfers := syncMapSwapSlice(t.ClientConMap)
	for _, v := range transfers {
		t.CloseClientConn(v)
	}
	// 再关闭监听对象
	if err = t.listener.Close(); err != nil {
		logger.Error("close server: %v error,errMsg: %v", t.Addr.String(), err)
		return err
	}
	logger.Info("TCPServer listener close finish")
	// 将监听对象置空
	t.listener = nil
	return
}

// 关闭指定客户端
func (t *TCPServer) CloseClientConn(client *TCPTransfer) {
	logger.Info("TCP CloseClientConn......")
	if t == nil || client == nil {
		return
	}
	logger.Info("before delete clientInfo:%+v", connMapPrint(t.ClientConMap))
	t.ClientConMap.Delete(client.ConnKey)
	logger.Info("after delete clientInfo:%+v", connMapPrint(t.ClientConMap))
	client.closeClientConn()
}

// 关闭客户端连接
func (t *TCPTransfer) closeClientConn() (err error) {
	logger.Info("TCPServer close clientConn , clientAddr:%v", t.Conn.RemoteAddr().String())
	if t.ConnState {
		if err := t.Conn.Close(); err != nil {
			logger.Error("server close client conn: %v error,errMsg: %v", t.Conn.LocalAddr().String(), err)
		}
		t.ConnState = false
	}
	return
}

// 接收监听并处理客户端handle
// 1. 获取客户端监听
// 2. 判断客户端连接数量是否超过限制值，超过则强行关闭客户端连接，不处理该客户端请求
// 3. 将客户端连接纳入连接管理，同时开启连接超时监听
// 4. 捕获panic 异常， 让TCP服务仍然可用
// 5. 利用 goroutine 去执行客户端定义的handle
func (t *TCPServer) acceptAndHandle() (err error) {
	var conn net.Conn
	if t.listener == nil {
		return errors.New("tcp: server listener is nil")
	}

	// 1. 处理客户端请求：
	if conn, err = t.listener.Accept(); err != nil {
		logger.Error("tcp: server( %s ) listener accept failed, errMsg: %s", t.Addr.String(), err.Error())
		// 直接断掉server服务 不执行下一轮的
		return
	}

	// 2. 控制客户端连接数量
	if concurrentUtil.SyncMapSize(&t.ClientConMap) >= t.ClientMaxNum {
		logger.Debug("refuse tcp client conn, curClinetMap:%v", connMapPrint(t.ClientConMap))
		logger.Error("tcp: curTCPService not support client conn, clientMaxNum: %d, now has clientNum: %d", t.ClientMaxNum, concurrentUtil.SyncMapSize(&t.ClientConMap))
		conn.Write([]byte("client over maxConnNum"))
		conn.Close()
		return
	}
	// 3.将客户端放入 map 中
	clientCon := NewTCPTransfer(conn)
	t.ClientConMap.Store(clientCon.ConnKey, clientCon)
	logger.Info("tcp: clientInfo -> `%#v`", t.ClientConMap)
	//开启客户端超时处理逻辑：
	if !t.clientConnTimeOutState {
		t.clientConnTimeOutState = true
		go func(t *TCPServer) {
			logger.Info("open tcp server listener client conn timeout check.....,tcpServer:%v", t.Addr)
			t.closeConnTimeOut()
		}(t)
	}

	//4. 处理客户端逻辑
	go func(t *TCPServer, cliConn *TCPTransfer) {
		logger.Info("================ TCP SERVER( %s <- %s ) START ================", t.Addr.String(), cliConn.ConnKey)
		// 1. 捕获潜在的异常 是否会退出上层的for循环是由PanicHandler决定的 默认情况不会退出for循环
		defer func() {
			if i := recover(); i != nil {
				if t.PanicHandler == nil {
					defaultPanicHandler(i)
				} else {
					t.PanicHandler(i)
				}
				t.CloseClientConn(cliConn)
			}
			logger.Info("================ TCP SERVER( %s <- %s ) END ================", t.Addr.String(), cliConn.ConnKey)
		}()

		// 用于初始化
		if t.InitHandler != nil {
			logger.Info("tcp: server( %s <- %s ) exec InitHandler...", t.Addr.String(), cliConn.ConnKey)
			if err = t.InitHandler(cliConn); err != nil {
				logger.Info("================ TCP SERVER( %s <- %s ) INITHANDLER FAILED ================", t.Addr.String(), cliConn.ConnKey)
				t.CloseClientConn(cliConn)
				return
			}
		}
		// 执行正常流程
		if t.Handler != nil {
			logger.Info("tcp: server( %s <- %s ) exec Handler...", t.Addr.String(), cliConn.ConnKey)
			t.Handler(cliConn)
		}
	}(t, clientCon)

	return
}

// 默认是捕获panic但不会因为有panic而退出更上层的for循环
func defaultPanicHandler(ifPanic interface{}) {
	if ifPanic != nil {
		logger.Error("================ TCP PANIC, errMsg: `%v` ================", ifPanic)
	}
}

// 刷新TCP 服务中的客户端连接活跃时间
func (t *TCPServer) FlushTcpClientActiveTime(connkey string) {
	if t == nil {
		return
	}
	if clientTransfer, ok := t.ClientConMap.Load(connkey); ok {
		logger.Debug("FlushTcpClientActiveTime, client:%v", connkey)
		transfer := clientTransfer.(*TCPTransfer) //因为是指针类型，所以直接可以修改值
		transfer.ClientLastActiveTime = time.Now().UnixMilli()
	}
}

// 关闭tcp 服务中客户端超时的连接
func (t *TCPServer) closeConnTimeOut() {
	count := 0
	for {
		time.Sleep(time.Millisecond * 500)
		if t == nil {
			continue
		}
		count++
		// 遍历当前ClientConMap
		transfers := syncMapSwapSlice(t.ClientConMap)
		if count%6 == 0 {
			logger.Debug("interval 3 seconds, clientMap:%v", connMapPrint(t.ClientConMap))
		}
		for _, v := range transfers {
			if t.clientConnTimeOut == -1 {
				continue
			}
			if v.ClientLastActiveTime+t.clientConnTimeOut < time.Now().UnixMilli() {
				logger.Info("tcp Server close timeout client %v ,lastActiveTime:%v, hasExistTime: %v(ms)", v.ConnKey, v.ClientLastActiveTime, time.Now().UnixMilli()-v.ClientLastActiveTime)
				t.CloseClientConn(v)
			}
		}
	}
}

// 遍历syncMap并转化为切片类型的 *TCPTransfer
func syncMapSwapSlice(m sync.Map) (transfers []*TCPTransfer) {
	m.Range(func(key, value interface{}) bool {
		v := value.(*TCPTransfer)
		transfers = append(transfers, v)
		return true
	})
	return
}

// connMapPrint 打印server的ClientConMap
func connMapPrint(m sync.Map) string {
	var (
		transfers   = syncMapSwapSlice(m)
		transferStr strings.Builder
	)
	transferStr.WriteString(connMapPrintPrefix)
	for i, tr := range transfers {
		fmt.Fprintf(&transferStr, transferStrFmt, i, tr.ConnKey, tr.ConnState, tr.ClientLastActiveTime, time.Now().UnixMilli()-tr.ClientLastActiveTime)
	}
	transferStr.WriteString(connMapPrintSuffix)
	return transferStr.String()
}
