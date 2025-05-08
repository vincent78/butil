package netUtil

import (
	"github.com/gin-gonic/gin"
	"github.com/vincent78/butil/logger/logger1"
	"github.com/vincent78/butil/net/gcgin"
)

const (
	HTTP = "http://"
)

type HttpServer struct {
	Addr string
	Gin  *gin.Engine
}

// 创建HttpServer
func NewHttpServer(addr string) (h *HttpServer) {
	logger1.Info("new HttpServer , open addr: [%v]", addr)
	h = &HttpServer{
		Addr: addr,
		Gin:  gcgin.InitEngine(true),
	}
	return h
}

// 添加处理Handle
func (httpServer *HttpServer) AddHandle(group string, path string, handle func(c *gin.Context)) *HttpServer {
	logger1.Info("AddHandle group: [%v] router : [%v] ", group, path)
	httpServer.Gin.Group(group).Any(path, handle)
	return httpServer
}

// 批量添加处理Handle
func (httpServer *HttpServer) AddHandles(group string, handlers map[string]func(c *gin.Context)) *HttpServer {
	logger1.Info("AddHandles group: [%v] handle.routers: [%v]", group, handlers)
	routerGroup := httpServer.Gin.Group(group)
	for k, v := range handlers {
		routerGroup.Any(k, v)
	}
	return httpServer
}

// 启动HttpServer
func (httpServer *HttpServer) Run() {
	logger1.Info("Run HttpServer , addr: [%v]", httpServer.Addr)
	gcgin.StartServer(httpServer.Addr, nil, httpServer.Gin)
}

// 关闭HttpServer
func (httpServer *HttpServer) Close() {
	if httpServer == nil {
		return
	}
	logger1.Info("Close HttpServer , addr: [%v]", httpServer.Addr)
	listener := gcgin.GetListener(httpServer.Addr)
	if listener != nil {
		(*listener).Close()
	}
	gcgin.RemoveListener(httpServer.Addr)
	logger1.Info("Close HttpServer , addr: [%v] finish", httpServer.Addr)
}
