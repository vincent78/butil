package server

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/logger/logger1"
	"github.com/vincent78/butil/net/gcws/common"
	"github.com/vincent78/butil/timewheel"
	"github.com/vincent78/butil/token"
	"github.com/vincent78/butil/utils/strUtil"
	"github.com/vincent78/butil/utils/timeUtil"
	"net/http"
	"time"
)

// WSServer is a middleman between the websocket connection and the Hub.
type WSServer struct {
	ctx           context.Context    // context
	cancelFunc    context.CancelFunc // cancelFunc
	conn          *websocket.Conn    // 原始连接对象
	token         string             // http的token
	ID            string             // 用户身份的唯一标识
	RemoteIP      string             // 远程的IP
	connectedTime int64              // 登陆的时间
	lastRespTime  int64              // 最后收到回复的时间

	writeCH chan []byte // 发送的内容  CMD & RESP
	readCH  chan []byte // 接收的内容  CMD & RESP
	//subscribeList []*WSServerSubscriber // 所有的监听者
	cmdHandlers map[string]WSServerCmdHandlerInterface
	config      *WSServerConfig // 配置文件

	pingTask *timewheel.Task // 时间轮中的心跳任务
}

func NewWSServer(conn *websocket.Conn, tk string, conf *WSServerConfig) *WSServer {
	server := &WSServer{
		conn:          conn,
		token:         tk,
		ID:            "",
		RemoteIP:      conn.RemoteAddr().String(),
		connectedTime: timeUtil.NowMillis(),
		writeCH:       make(chan []byte, maxMessagePool),
		readCH:        make(chan []byte, maxMessagePool),
		cmdHandlers:   make(map[string]WSServerCmdHandlerInterface, 0),
		config:        conf,
	}
	server.ctx, server.cancelFunc = context.WithCancel(context.Background())
	return server
}

func (c *WSServer) Resp(resp *common.WSCmd) {
	c.writeCH <- resp.Bytes()
}

func (c *WSServer) Send(cmd *common.WSCmd) {
	if cmd.ID == "" {
		cmd.ID = token.NewShort()
	}
	c.writeCH <- cmd.Bytes()
}

func (c *WSServer) Dispose() {
	defer func() {
		if r := recover(); r != nil {
			logger1.ErrorByName(global.LogFileWSSName, c.NormalLogger("dispose with error: %v", r))
		}
	}()
	if c.cancelFunc != nil {
		c.cancelFunc()
	}
	if c.pingTask != nil {
		global.Timewheel.Remove(c.pingTask)
	}
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	close(c.writeCH)
}

// readPump pumps messages from the websocket connection to the Hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *WSServer) readPump() {
	logger1.InfoByName(global.LogFileWSSName, c.NormalLogger("readPump start"))
	defer func() {
		logger1.InfoByName(global.LogFileWSSName, c.NormalLogger("readPump end"))
	}()

	defer func() {
		if r := recover(); r != nil {
			logger1.ErrorByName(global.LogFileWSSName, c.NormalLogger("readPump: %v", r))
		}
	}()

	c.conn.SetReadLimit(maxMessageSize)
	for {
		mType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseNoStatusReceived,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				logger1.ErrorByName(global.LogFileWSSName, "error: %v", err)
				continue
			}
			HubManager.unregister <- c
			return
		}

		c.lastRespTime = timeUtil.NowMillis()
		if mType != websocket.TextMessage {
			logger1.WarnByName(global.LogFileWSSName, c.NormalLogger("receive msg type %v - %v", mType, strUtil.Bytes2String(message)))
			continue
		} else {
			logger1.DebugByName(global.LogFileWSSName, c.InLogger("%v", strUtil.Bytes2String(message)))
			if cmd, err := common.ParseCmd(message); err == nil {
				if cmd.Cmd == "" && cmd.PID != "" {
					logger1.DebugByName(global.LogFileWSSName, c.NormalLogger("receive the resp without handler"))
				} else {
					key := cmd.HandlerKey()
					if handler, exist := c.cmdHandlers[key]; exist && handler != nil {
						handler.DoAction(c, cmd)
					} else {
						logger1.ErrorByName(global.LogFileWSSName, c.NormalLogger("the cmd[%v] handler is not exist or not current", cmd.HandlerKey()))
					}
				}
			} else {
				logger1.ErrorByName(global.LogFileWSSName, c.NormalLogger("cmd parse error: %v ", err.ToString()))
				if cmd.Cmd != "" {
					c.Send(cmd.FailureByErrModel(err))
				}
			}
		}
	}
}

// writePump pumps messages from the Hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *WSServer) writePump() {
	logger1.InfoByName(global.LogFileWSSName, c.NormalLogger("writePump start"))
	defer func() {
		logger1.InfoByName(global.LogFileWSSName, c.NormalLogger("writePump end"))
	}()

	defer func() {
		if r := recover(); r != nil {
			logger1.ErrorByName(global.LogFileWSSName, c.NormalLogger("writePump: %v", r))
		}
	}()

	for {
		select {
		case message, _ := <-c.writeCH:
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				logger1.ErrorByName(global.LogFileWSSName, c.NormalLogger(err.Error()))
			} else {
				logger1.DebugByName(global.LogFileWSSName, c.OutLogger(strUtil.Bytes2String(message)))
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request, f func(s *WSServer)) {
	tk := r.Header.Get(global.RequestTokenKey)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger1.ErrorByName(global.LogFileWSSName, err.Error())
		return
	}

	if "" == tk {
		tk = token.NewShort()
	}
	r.Header.Set(global.RequestTokenKey, tk)

	server := NewWSServer(conn, tk, HubManager.WSServerConfig)
	server.initCmdHandler()

	HubManager.register <- server
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go server.writePump()
	go server.readPump()

	conn.SetCloseHandler(func(code int, text string) error {
		logger1.InfoByName(global.LogFileWSSName, server.NormalLogger("handle close action: %v - %v", code, text))
		HubManager.unregister <- server
		return nil
	})

	if server.config.PingInterval > 0 {
		server.pingTask = global.Timewheel.AddCron(time.Duration(server.config.PingInterval)*time.Second, server.ping)
	}

	if f != nil {
		f(server)
	}
}

func (c *WSServer) ping() {
	cmd := common.NewWSCmdByCategory(common.CmdCategorySys, common.CmdNamePing)
	cmd.Data = timeUtil.NowStr()
	c.Send(cmd)
}

func (c *WSServer) RegisterCmdHandler(category, cmd string, f WSServerCmdHandlerInterface) {
	key := fmt.Sprintf("%v:%v", category, cmd)
	c.cmdHandlers[key] = f
	//c.cmdHandlers[fmt.Sprintf("%v:%v", category, cmd)] = f
}
