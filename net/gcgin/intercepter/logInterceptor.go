package intercepter

import (
	"bytes"
	"fmt"
	"gitee.com/vincent78/gcutil/global"
	"gitee.com/vincent78/gcutil/logger"
	"gitee.com/vincent78/gcutil/utils/mapUtil"
	"gitee.com/vincent78/gcutil/utils/strUtil"
	"github.com/gin-gonic/gin"
	"io"
	"regexp"
	"strings"
	"time"
)

const MaxPrintBodyLen = 512

var showRequetHeader = false
var showResponseHeader = false
var showResponseBody = true

type bodyLogWriter struct {
	gin.ResponseWriter
	bodyBuf *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	//memory copy here!
	w.bodyBuf.Write(b)
	return w.ResponseWriter.Write(b)
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		token := c.Request.Context().Value(global.RequestTokenKey)
		blackPath := false
		for k, _ := range global.GinRepBodyBlackPath {
			if ok, _ := regexp.Match(k, []byte(path)); ok {
				blackPath = true
				break
			}
		}

		var blw bodyLogWriter
		if showResponseBody && !blackPath {
			blw = bodyLogWriter{bodyBuf: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw
		}

		data := GenerateRequestBody(c)

		obj := make(map[string]interface{})
		obj["method"] = c.Request.Method
		obj["path"] = path
		if len(query) > 0 {
			obj["query"] = query
		}

		if data != "" {
			obj["body"] = data
		}

		if c.Request.Form != nil {
			obj["form"] = c.Request.Form
		}

		if showRequetHeader && c.Request.Header != nil {
			obj["header"] = c.Request.Header
		}

		logger.InfoByName(global.LogFileHttpName, fmt.Sprintf("---> [%v][%v] :%v", GetClientIp(c), token, mapUtil.ToJsonStr(obj)))

		c.Next()
		cost := time.Since(start)

		respLogMap := make(map[string]interface{})
		respLogMap["cost"] = cost
		if c.Errors != nil {
			respLogMap["errors"] = c.Errors.ByType(gin.ErrorTypePrivate).String()
		}

		if showResponseBody && !blackPath {
			respLogMap["body"] = blw.bodyBuf.String()
		} else {
			respLogMap["body"] = "the path in black list,so there show nothing"
		}
		logger.InfoByName(global.LogFileHttpName, fmt.Sprintf("<--- [%v][%v] :%v", GetClientIp(c), token, mapUtil.ToJsonStr(respLogMap)))
	}
}

func GenerateRequestBody(c *gin.Context) string {
	bts, err := c.GetRawData()
	if err != nil {
		logger.DebugByName(global.LogFileHttpName, "failed to get request body")
		return ""
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bts))
	body := strUtil.Bytes2String(bts)
	body = strings.ReplaceAll(body, "\n", "")
	return body
}

func GetClientIp(g *gin.Context) string {

	reqIP := g.ClientIP()
	if emptyIp(reqIP) {
		// 注意，这里要求nginx中必须配置 X-Real-Ip
		reqIP = g.Request.Header.Get("X-Real-Ip")
	}
	if emptyIp(reqIP) {
		// 注意，这里要求nginx中必须配置 X-Real-Ip
		reqIP = g.Request.Header.Get("X-real-ip")
	}

	if emptyIp(reqIP) {
		reqIP = g.Request.Header.Get("X-Forward-For")
	}

	if emptyIp(reqIP) {
		reqIP = "127.0.0.1"
	}
	return reqIP
}

func emptyIp(ip string) bool {
	return ip == "" || strings.Contains(ip, "127.0.0.1") || ip == "::1"
}
