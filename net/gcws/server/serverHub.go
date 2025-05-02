package server

import (
	"context"
	"gitee.com/vincent78/gcutil/global"
	"gitee.com/vincent78/gcutil/logger"
	"gitee.com/vincent78/gcutil/net/gcws/common"
	"github.com/gin-gonic/gin"
)

// Hub maintains the set of active servers and broadcasts messages to the
// servers.
type Hub struct {
	// Registered servers.
	servers map[*WSServer]struct{}
	// login server
	loginServers map[string]*WSServer
	// Inbound messages from the servers.
	broadcast chan []byte
	// register requests from the servers.
	register chan *WSServer
	// unregister requests from servers.
	unregister chan *WSServer

	WSServerConfig *WSServerConfig
	context        context.Context
	cancelFunc     context.CancelFunc
}

func InitWSByGin(e *gin.Engine, path string, f func(s *WSServer), conf *WSServerConfig) {
	if HubManager == nil {
		HubManager = newHub()
		go HubManager.run()
		common.InitTimewheel(global.LogFileWSSName)
	}
	HubManager.WSServerConfig = conf
	e.GET(path, func(c *gin.Context) {
		c.Request.Header["Origin"] = nil
		serveWs(c.Writer, c.Request, f)
	})
	if len(conf.InfoPath) > 0 {
		e.GET(conf.InfoPath, wsInfo)
	}
}

func newHub() *Hub {
	ctx, f := context.WithCancel(context.Background())
	return &Hub{
		context:      ctx,
		cancelFunc:   f,
		broadcast:    make(chan []byte, 10),
		register:     make(chan *WSServer, 10),
		unregister:   make(chan *WSServer, 10),
		servers:      map[*WSServer]struct{}{},
		loginServers: map[string]*WSServer{},
	}
}

func (h *Hub) run() {
	defer func() {
		if r := recover(); r != nil {
			logger.ErrorByName(global.LogFileWSSName, "%v", r)
		}
	}()
	for {
		select {
		// 客户端连接
		case server := <-h.register:
			h.servers[server] = struct{}{}
			logger.InfoByName(global.LogFileWSSName, server.NormalLogger("connected"))
		// 客户端断开连接
		case server := <-h.unregister:
			if server.ID != "" {
				h.Logout(server)
			}
			if _, ok := h.servers[server]; ok {
				server.Dispose()
				delete(h.servers, server)
				logger.InfoByName(global.LogFileWSSName, server.NormalLogger("disconnected"))
			} else {
				logger.WarnByName(global.LogFileWSSName, server.NormalLogger("can't find anonymity"))
			}
		case bytes := <-h.broadcast:
			for server := range h.servers {
				server.writeCH <- bytes
			}
			for _, server := range h.loginServers {
				server.writeCH <- bytes
			}
		case <-h.context.Done():
			h.Close()
		}
	}
}

func (h *Hub) Login(server *WSServer) {
	logger.InfoByName(global.LogFileWSSName, server.NormalLogger("login"))
	h.loginServers[server.ID] = server
	delete(h.servers, server)
}

func (h *Hub) Logout(server *WSServer) {
	logger.InfoByName(global.LogFileWSSName, server.NormalLogger("logout"))
	delete(h.loginServers, server.ID)
	server.ID = ""
	h.servers[server] = struct{}{}
}

func (h *Hub) UserCount() int {
	return len(h.loginServers)
}

func (h *Hub) AnonymityCount() int {
	return len(h.servers)
}

func (h *Hub) CloseServer(c *WSServer) {
	if c != nil {
		h.unregister <- c
	} else {
		logger.ErrorByName(global.LogFileWSSName, "the server is nil")
	}
}

func (h *Hub) Close() {
	for server, _ := range h.servers {
		h.CloseServer(server)
	}

	for _, server := range h.loginServers {
		h.CloseServer(server)
	}
}
