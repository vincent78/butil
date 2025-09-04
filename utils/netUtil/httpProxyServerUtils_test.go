package netUtil

import (
	"fmt"
	io "io"
	"net/http"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/logger/logger1"
	"github.com/vincent78/butil/utils/strUtil"
)

func initlogger() {
	path := "/tmp/net"
	c1 := logger1.NewLogConfig()
	c1.Path = path
	logger1.NewLogger(c1)
	c2 := logger1.NewLogConfig()
	c2.Path = path
	c2.Name = global.LogFileNetName
	logger1.NewLogger(c2)
}

func TestHttpServer(t *testing.T) {
	initlogger()
	httpServer := NewHttpServer(":9001")
	httpServer.AddHandle("/api", "/*a", proxyHandler).Run()
	defer httpServer.Close()
}

func proxyHandler1(c *gin.Context) {
	accepted := c.Accepted
	logger1.InfoByName(global.LogFileNetName, "accepted :%v", accepted)
	fmt.Printf("gin.Context:%+v", c)
	uri := c.Request.RequestURI
	fmt.Println(uri)

	header := c.Request.Header
	for k, v := range header {
		logger1.InfoByName(global.LogFileNetName, "k: %v, v:%v", k, v)
	}

	questType := c.Request.Method
	logger1.InfoByName(global.LogFileNetName, "requestType: %v", questType)

	rcsUrl := "http://172.31.239.56:9001"

	resp, err := http.Post(rcsUrl+uri, c.Request.Header.Get("Content-Type"), c.Request.Body)

	if err != nil {
		logger1.InfoByName(global.LogFileNetName, "err:%v", err)
	}
	all, _ := io.ReadAll(resp.Body)
	c.Writer.Write(all)

}

func proxyHandler(c *gin.Context) {

	rcsUrl := "http://172.31.239.56:9001"

	requestType := c.Request.Method // 获取请求方式
	targetUrl := rcsUrl + c.Request.RequestURI
	c.Writer.Header().Add("Content-Type", "application/json")
	//转发post 请求
	if requestType == "POST" {
		logger1.Info("forward post request , targetUrl: %v", targetUrl)
		contentLength, _ := strconv.Atoi(c.Request.Header.Get("Content-Length"))
		var buf = make([]byte, contentLength, contentLength)
		read, _ := c.Request.Body.Read(buf[:contentLength])
		receiverMes := string(buf[:read])
		logger1.Info("request para:%v", receiverMes)
		var headerMap = make(map[string]string)
		for k, v := range c.Request.Header {
			headerMap[k] = v[0]
		}
		// 发送POST 请求
		//result := gchttp.PostRequest(targetUrl, headerMap, receiverMes)
		// 处理成功信息
		//if result.Code == model.Success {
		//	str := strUtil.ToJsonStr(result.Data)
		//	logger1.Info("receive requestUrl: %v response str:%v", targetUrl, str)
		//	c.Writer.Write(strUtil.String2Bytes(str))
		//	return
		//}
		// 处理失败信息
		//writeErrorResponse(targetUrl, c.Writer, result.Message)
		return
	}
	//转发GET 请求
	if requestType == "GET" {
		// 发送GET 请求
		resp, err := http.Get(targetUrl)
		//处理失败响应信息
		if err != nil {
			writeErrorResponse(targetUrl, c.Writer, err.Error())
			return
		}
		//处理成功响应信息
		writeSuccessResponse(targetUrl, c.Writer, resp)
		return
	}
}

// 将真实服务返回的response 结果写入给真实的请求客户端
func writeSuccessResponse(targetUrl string, relHttpResp gin.ResponseWriter, proxyHttpResp *http.Response) {
	readAll, err := io.ReadAll(proxyHttpResp.Body)
	if err != nil {
		logger1.Error("read forward post request: %v response error,errInfo:%v ", targetUrl, err)
		return
	}
	logger1.Info("receive requestUrl: %v response result:%v", targetUrl, string(readAll))
	relHttpResp.Write(readAll)
}

// 将代理Http服务出现的异常进行错误响应回复
func writeErrorResponse(targetUrl string, relHttpResp gin.ResponseWriter, errMes string) {
	logger1.Error("receive requestUrl: %v request failed!  errInfo:%v", targetUrl, errMes)
	response := make(map[string]string)
	response["success"] = "false"
	response["code"] = "QSH000000"
	response["message"] = errMes
	resStr := strUtil.ToJsonStr(response)
	logger1.Error("receive requestUrl: %v request failed!  result:%v", targetUrl, resStr)
	relHttpResp.Write([]byte(resStr))
}
