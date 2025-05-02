package netUtil

import (
	"gitee.com/vincent78/gcutil/logger"
	"gitee.com/vincent78/gcutil/net/gcgin"
	"github.com/gin-gonic/gin"
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
	logger.Info("new HttpServer , open addr: [%v]", addr)
	h = &HttpServer{
		Addr: addr,
		Gin:  gcgin.InitEngine(true),
	}
	return h
}

// 添加处理Handle
func (httpServer *HttpServer) AddHandle(group string, path string, handle func(c *gin.Context)) *HttpServer {
	logger.Info("AddHandle group: [%v] router : [%v] ", group, path)
	httpServer.Gin.Group(group).Any(path, handle)
	return httpServer
}

// 批量添加处理Handle
func (httpServer *HttpServer) AddHandles(group string, handlers map[string]func(c *gin.Context)) *HttpServer {
	logger.Info("AddHandles group: [%v] handle.routers: [%v]", group, handlers)
	routerGroup := httpServer.Gin.Group(group)
	for k, v := range handlers {
		routerGroup.Any(k, v)
	}
	return httpServer
}

// 启动HttpServer
func (httpServer *HttpServer) Run() {
	logger.Info("Run HttpServer , addr: [%v]", httpServer.Addr)
	gcgin.StartServer(httpServer.Addr, nil, httpServer.Gin)
}

// 关闭HttpServer
func (httpServer *HttpServer) Close() {
	if httpServer == nil {
		return
	}
	logger.Info("Close HttpServer , addr: [%v]", httpServer.Addr)
	listener := gcgin.GetListener(httpServer.Addr)
	if listener != nil {
		(*listener).Close()
	}
	gcgin.RemoveListener(httpServer.Addr)
	logger.Info("Close HttpServer , addr: [%v] finish", httpServer.Addr)
}
