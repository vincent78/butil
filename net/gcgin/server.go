package gcgin

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/logger/logger1"
	itcp "github.com/vincent78/butil/net/gcgin/intercepter"
	"github.com/vincent78/butil/net/gcgin/model"
	"github.com/vincent78/butil/sys"
	"github.com/vincent78/butil/token"
	"github.com/vincent78/butil/utils/timeUtil"
)

// http 监听map,用于外面获取监听对象
var listenerMap sync.Map

func InitEngine(debug bool) *gin.Engine {
	e := gin.New()
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	config := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "UPDATE", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		//AllowHeaders:    []string{"Content-Length", "Authorization", "Content-Type", "Origin", "x-access-token"},
		AllowHeaders: []string{"*"},
		//ExposeHeaders:    []string{"Content-Length", "Authorization", "Content-Type", "Origin", "x-access-token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	logger1.DebugByName(global.LogFileHttpName, "gin config")

	e.Use(cors.New(config))
	e.Use(itcp.GinRecovery())
	e.Use(itcp.TokenInterceptor())
	e.Use(itcp.GinLogger())

	se := e.Group("/sys")
	se.GET("timestamp", timestampHandler)
	se.GET("token", tokenHandler)
	se.GET("version", versionHandler)
	//if debug {
	//	caller, file, line, ok := runtime.Caller(1)
	//	logger.DebugByName(global.LogFileHttpName, "caller:%v, file:%v, line:%v, ok:%v", caller, file, line, ok)
	//}
	return e
}

func SyncStartServer(addr string, g *gin.Engine) {
	// start the web server
	c := make(chan struct{})
	go StartServer(addr, c, g)
	<-c
}

func StartServer(addr string, c chan struct{}, g *gin.Engine) {
	// “tcp”、“tcp4”、“tcp6”、"unix "或 “unixpacket”
	// 强制绑定在ipv4
	l, _ := net.Listen("tcp4", addr)
	if c != nil {
		c <- struct{}{}
	}
	listenerMap.Store(addr, l)
	g.HandleMethodNotAllowed = true
	g.NoRoute(recover400)
	g.Use(recover500)
	logger1.InfoByName(global.LogFileHttpName, "------------------------------------------")
	logger1.InfoByName(global.LogFileHttpName, "web port:\t\t %v", addr)
	logger1.InfoByName(global.LogFileHttpName, "current pid: %v", os.Getpid())
	logger1.InfoByName(global.LogFileHttpName, "------------------------------------------")

	err := g.RunListener(l)
	if err != nil {
		logger1.ErrorByName(global.LogFileHttpName, "start server error: %v", err)
	}
}

// GetListener 获取http 监听对象
func GetListener(addr string) *net.Listener {
	if load, ok := listenerMap.Load(addr); ok {
		return load.(*net.Listener)
	}

	return nil
}

// RemoveListener 移除http 监听对象
func RemoveListener(addr string) {
	listenerMap.Delete(addr)
}

func timestampHandler(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("%d", timeUtil.NowMillis()))
}

func tokenHandler(c *gin.Context) {
	c.String(http.StatusOK, token.NextStrToken())
}

func versionHandler(c *gin.Context) {
	obj := sys.Version()
	c.JSON(http.StatusOK, obj)
}

func recover400(c *gin.Context) {
	c.JSON(http.StatusOK, model.ErrorHttp404(fmt.Sprintf("%v: %v", c.Request.Method, c.Request.RequestURI)))
}

func recover500(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger1.ErrorByName(global.LogFileHttpName, "panic: %v\n", r)
			//debug.PrintStack()
			c.JSON(200, gin.H{
				"code":    500,
				"message": "服务器内部错误",
			})
		}
	}()
	c.Next()
}
