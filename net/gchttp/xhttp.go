/* *
 * http Resty V2
 */

package gchttp

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var once sync.Once
var xhttpClient *xhttp

type xhttp struct {
	ctx context.Context
	*resty.Client
}

// NewXHTTP 创建 *xhttp
func NewXHTTP(ctx context.Context) *xhttp {
	once.Do(func() {
		xhttpClient = &xhttp{
			ctx,
			resty.NewWithClient(&http.Client{
				Transport: &http.Transport{
					MaxIdleConns:        500,              // 全局空闲连接
					MaxIdleConnsPerHost: 100,              // 单个host空闲连接
					MaxConnsPerHost:     200,              // 单个host最大连接
					IdleConnTimeout:     90 * time.Second, // 空闲超时

					TLSHandshakeTimeout:   10 * time.Second, // TLS 握手超时
					ExpectContinueTimeout: 1 * time.Second,  // 100-Continue 超时
					ResponseHeaderTimeout: 60 * time.Second, // 等待响应头超时

					DialContext: (&net.Dialer{
						Timeout:   10 * time.Second, // TCP 建连超时
						KeepAlive: 30 * time.Second, // TCP keepalive 心跳
					}).DialContext,

					// 跳过证书验证
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
				Timeout: 120 * time.Second, // 整个请求超时（含建连+读写）
			}),
		}

		runtime.SetFinalizer(xhttpClient, func(h *xhttp) {
			h.close()
		})
	})

	return xhttpClient
}

// 关闭http底层tcp链接
func (h *xhttp) close() {
	runtime.KeepAlive(xhttpClient)
	if transport, err := h.Transport(); err == nil {
		transport.CloseIdleConnections()
	}
}

// Close 程序退出时,关闭http底层tcp链接
func Close() {
	if xhttpClient != nil {
		xhttpClient.close()
	}
}
