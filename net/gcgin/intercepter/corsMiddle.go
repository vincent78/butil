package intercepter

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查是否允许的源
		if isAllowedOrigin(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// 设置其他 CORS 头
		c.Header("Access-Control-Allow-Methods", "*") // 允许所有方法
		// c.Header("Access-Control-Allow-Headers", "Authorization, traceparent, tracestate, cache-control")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// isAllowedOrigin 检查源是否被允许
func isAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}

	// 移除协议前缀，获取主机部分
	var host string
	if strings.HasPrefix(origin, "http://") {
		host = strings.TrimPrefix(origin, "http://")
	} else if strings.HasPrefix(origin, "https://") {
		host = strings.TrimPrefix(origin, "https://")
	} else {
		return false
	}

	// 分离主机和端口
	hostParts := strings.Split(host, ":")
	hostname := hostParts[0]

	// 1. 检查 localhost:port
	if hostname == "localhost" {
		return true
	}

	// 2. 检查 127.0.0.1:port
	if hostname == "127.0.0.1" {
		return true
	}

	// 3. 检查 192.168.*:port
	matched, _ := regexp.MatchString(`^192\.168\.\d+\.\d+$`, hostname)
	if matched {
		return true
	}

	// 4. 检查 *.helix.city
	if strings.HasSuffix(hostname, ".helix.city") {
		return true
	}

	// 5. 检查 *.justraise.pro
	if strings.HasSuffix(hostname, ".justraise.pro") {
		return true
	}

	return false
}
