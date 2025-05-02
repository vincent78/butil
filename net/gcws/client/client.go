package client

import (
	"bytes"
	"context"
	"github.com/gorilla/websocket"
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/logger"
	"github.com/vincent78/butil/net/gcws/common"
	"github.com/vincent78/butil/timewheel"
	"github.com/vincent78/butil/token"
	"github.com/vincent78/butil/utils/strUtil"
	"github.com/vincent78/butil/utils/timeUtil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type WSClient struct {
	url          *url.URL
	conn         *websocket.Conn
	ClientAddr   string // 本地监听的地址与端口  0.0.0.0:10023
	writeChan    chan []byte
	readChan     chan []byte
	ctx          context.Context
	cancelFunc   context.CancelFunc
	successTime  int64
	pingTask     *timewheel.Task
	lastRespTime int64
	cmdHandlers  map[string]WSClientCmdHandlerInterface
	config       *WSClientConfig
}

func NewWSClient(schema, addr, path string) *WSClient {
	return NewWSClientByConfig(schema, addr, path, NewWSClientConfig())
}

func NewWSClientByConfig(schema, addr, path string, conf *WSClientConfig) *WSClient {
	common.InitTimewheel(global.LogFileWSCName)
	ctx, cancelFunc := context.WithCancel(context.Background())
	client := &WSClient{
		url:         &url.URL{Scheme: schema, Host: addr, Path: path},
		conn:        nil,
		writeChan:   make(chan []byte, maxMessagePool),
		readChan:    make(chan []byte, maxMessagePool),
		ctx:         ctx,
		cancelFunc:  cancelFunc,
		ClientAddr:  "",
		successTime: timeUtil.NowMillis(),
		cmdHandlers: make(map[string]WSClientCmdHandlerInterface),
		config:      conf,
	}
	client.initSysCmdHander()
	return client
}

func (c *WSClient) Resp(resp *common.WSCmd) {
	c.writeChan <- resp.Bytes()
}

func (c *WSClient) Send(cmd *common.WSCmd) {
	if cmd.ID == "" {
		cmd.ID = token.NewShort()
	}
	//if _, exist := c.runningCmdMap[cmd.ID]; exist {
	//	cmd.ID = fmt.Sprintf("%v:%v", cmd.ID, token.NewShort())
	//}
	//
	//if _, exist := c.ignoreRecordMap[cmd.HandlerKey()]; !exist {
	//	c.runningCmdMap[cmd.ID] = cmd
	//}
	c.writeChan <- cmd.Bytes()
}

func (c *WSClient) Conn() error {
	common.InitTimewheel(global.LogFileWSCName)
	urlStr := c.url.String()
	logger.InfoByName(global.LogFileWSCName, "%v connecting to %s", c.LogPrefix(), urlStr)
	dialer := websocket.DefaultDialer
	if c.ClientAddr != "" {
		netDialer := &net.Dialer{}
		netDialer.LocalAddr, _ = net.ResolveTCPAddr("", c.ClientAddr)
		dialer.NetDial = func(network, addr string) (net.Conn, error) {
			return netDialer.DialContext(context.Background(), network, addr)
		}
	}

	var requestHeader http.Header
	var err error
	c.conn, err = c.reset(dialer, urlStr, requestHeader)
	if err == nil && c.conn != nil {

		c.ClientAddr = c.conn.LocalAddr().String()
		c.ClientAddr = c.ClientAddr[strings.LastIndex(c.ClientAddr, ":"):]

		c.conn.SetCloseHandler(func(code int, text string) error {
			logger.InfoByName(global.LogFileWSCName, "%v handle close event: %v  -  %v", c.LogPrefix(), code, text)
			return nil
		})

		go c.sendPump()
		go c.receivePump()
		if c.config.PingInterval > 0 {
			c.pingTask = global.Timewheel.AddCron(pingPeriod, c.ping)
		}
		logger.InfoByName(global.LogFileWSCName, "%v connected: %v", c.LogPrefix(), urlStr)
	}
	return err
}

func (c *WSClient) sendPump() {
	logger.DebugByName(global.LogFileWSCName, "%v: sendPump start", c.LogPrefix())
	defer func() {
		logger.DebugByName(global.LogFileWSCName, "%v: sendPump end", c.LogPrefix())
	}()

	defer func() {
		if r := recover(); r != nil {
			logger.Error("%v sendPump with error: %v", c.LogPrefix(), r)
		}
	}()
	for {
		select {
		case bs := <-c.writeChan:
			if bs != nil && len(bs) > 0 {
				logger.DebugByName(global.LogFileWSCName, "%v: %v", c.LogPrefixOut(), strUtil.Bytes2String(bs))
				err := c.conn.WriteMessage(websocket.TextMessage, bs)
				if err != nil {
					logger.ErrorByName(global.LogFileWSCName, "%v send error: %v", c.LogPrefix(), err)
				}
			}

		case <-c.ctx.Done():
			return
		}
	}
}

func (c *WSClient) receivePump() {
	logger.DebugByName(global.LogFileWSCName, "%v: receivePump start", c.LogPrefix())
	defer func() {
		logger.DebugByName(global.LogFileWSCName, "%v: receivePump end", c.LogPrefix())
	}()

	defer func() {
		if r := recover(); r != nil {
			logger.Error("%v receivePump with error: %v", c.LogPrefix(), r)
		}
	}()
	for {
		mType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseNoStatusReceived,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				logger.ErrorByName(global.LogFileWSCName, "error: %v", err)
				continue
			}
			if c.url != nil {
				global.Timewheel.AfterFunc(time.Second, func() {
					_, err := c.reConn()
					if err != nil {
						logger.ErrorByName(global.LogFileWSCName, "reconn error: %v", err)
					}
				})
			}
			return
		}

		if mType != websocket.TextMessage {
			logger.WarnByName(global.LogFileWSCName, "%v: receive the msg type - %v", mType)
			continue
		}

		message = bytes.TrimSpace(message)
		logger.DebugByName(global.LogFileWSCName, "%v: %v", c.LogPrefixIn(), strUtil.Bytes2String(message))
		c.lastRespTime = timeUtil.NowMillis()
		cmd, em := common.ParseCmd(message)
		if em != nil {
			logger.WarnByName(global.LogFileWSCName, "%v parseCmd[ %v ] error: %v", c.LogPrefix(), strUtil.Bytes2String(message), em)
			if cmd.PID == "" {
				c.Resp(cmd.FailureByErrModel(common.ErrorWSMParseCmd("")))
			}
		} else {
			DoCmdAction(c, cmd)
		}
	}
}

func (c *WSClient) reset(dialer *websocket.Dialer, urlstr string, reqHeader http.Header) (*websocket.Conn, error) {
	conn, _, err := dialer.Dial(urlstr, reqHeader)
	currNum := 0
	if err != nil {
		logger.WarnByName(global.LogFileWSCName, "%v connect error: %v ", c.LogPrefix(), err.Error())
		tick := time.NewTicker(time.Duration(config.ConnRetryInterval) * time.Second)
		for {
			select {
			case <-tick.C:
				if currNum < config.ConnMaxNum {
					logger.InfoByName(global.LogFileWSCName, "%v try connect num:[%v]", c.LogPrefix(), currNum)
					conn, _, err = dialer.Dial(urlstr, reqHeader)
				}
			}
			if err != nil {
				currNum = currNum + 1
				if currNum >= config.ConnMaxNum {
					logger.InfoByName(global.LogFileWSCName, "%v touch the try max num", c.LogPrefix())
					tick.Stop()
					break
				}
			} else {
				tick.Stop()
				break
			}
		}
	}
	return conn, err
}

func (c *WSClient) DisConn() {
	logger.InfoByName(global.LogFileWSCName, "%v: disconn ", c.LogPrefix())

	defer func() {
		if r := recover(); r != nil {
			logger.Error("%v disconn with error: %v", c.LogPrefix(), r)
		}
	}()
	c.url = nil // 设置为nil才不会自动重联

	if c.cancelFunc != nil {
		c.cancelFunc()
	}
	if c.conn != nil {
		c.conn.Close()
		close(c.readChan)
		close(c.writeChan)
		c.conn = nil
	}

	if c.pingTask != nil {
		if err := global.Timewheel.Remove(c.pingTask); err != nil {
			return
		}
	}
}

func (c *WSClient) reConn() (*WSClient, error) {
	logger.InfoByName(global.LogFileWSCName, "%v retry connect", c.LogPrefix())
	c.DisConn()
	client := NewWSClientByConfig(c.url.Scheme, c.url.Host, c.url.Path, config)
	err := client.Conn()
	return client, err
}

func (c *WSClient) SendCmd(cmd *common.WSCmd) {
	if cmd.Timestamp <= 0 {
		cmd.Timestamp = timeUtil.NowMillis()
	}
	c.writeChan <- cmd.Bytes()
}

func (c *WSClient) ping() {
	cmd := common.NewWSCmdByCategory(common.CmdCategorySys, common.CmdNamePing)
	cmd.Data = timeUtil.NowStr()
	c.Send(cmd)
}
