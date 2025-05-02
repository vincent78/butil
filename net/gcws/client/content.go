package client

import (
	"github.com/vincent78/butil/net/gcws/common"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// writeChan pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	maxMessagePool = 100
)

var DefaultServerPath = "/ws"

var defaultServerSchema = "ws"

var config *WSClientConfig

//var upgrader = websocket.Upgrader{
//	ReadBufferSize:  maxMessageSize,
//	WriteBufferSize: maxMessageSize,
//}
//
//var ignoreCmdMap = map[string]interface{}{
//	fmt.Sprintf("%v:%v", common.CmdCategorySys, common.CmdNamePing): struct{}{},
//	fmt.Sprintf("%v:%v", common.CmdCategorySys, common.CmdNamePong): struct{}{},
//}

type WSCmdClientHandler func(c *WSClient, cmd *common.WSCmd)
