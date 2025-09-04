package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vincent78/butil/net/gcws/common"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// writeCH pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024

	maxMessagePool = 100
)

var HubManager *Hub

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ignoreCmdMap = map[string]interface{}{
	fmt.Sprintf("%v:%v", common.CmdCategorySys, common.CmdNamePing): struct{}{},
	fmt.Sprintf("%v:%v", common.CmdCategorySys, common.CmdNamePong): struct{}{},
}
