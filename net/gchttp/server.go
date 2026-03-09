package gchttp

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/vincent78/butil/logger/logger1"
	"github.com/vincent78/butil/model"
	"github.com/vincent78/butil/token"
	gctime "github.com/vincent78/butil/utils/timeUtil"
)

type HttpHandler func(w http.ResponseWriter, r *http.Request)

var (
	httpRoutes = map[string]HttpHandler{
		"/timestamp": TimestampHandler,
		"/token01":   TokenHandler,
	}
)

func TimestampHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, gctime.NowMillis())
}

func TokenHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, token.NextStrToken())
}

func RegistHandle(name string, handler HttpHandler) {
	httpRoutes[name] = handler
}

func StartServer(addr string, c chan struct{}) {
	if addr == "" {
		logger1.Error("StartServer: the input params can't be nil.")
		return
	}
	for key, handler := range httpRoutes {
		logger1.Debug("gchttp register the handler: %v", key)
		http.HandleFunc(key, handler)
	}
	logger1.Debug("gchttp server : %v", addr)
	server := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	l, _ := net.Listen("tcp", addr)
	if c != nil {
		c <- struct{}{}
	}
	server.Serve(l)
	//server.ListenAndServe()
}

/*
*
openssl genrsa -out server.key 2048
openssl req -new -x509 -key server.key -out server.crt -days 365
*/
func StartTLSServer(addr, crt, key string) {
	if crt == "" || key == "" || addr == "" {
		logger1.Error("StartTLSServer: the input params can't be nil.")
		return
	}
	for key, handler := range httpRoutes {
		for key, handler := range httpRoutes {
			logger1.Debug("register the https handler: %v", key)
			http.HandleFunc(key, handler)
		}
		http.HandleFunc(key, handler)
	}
	logger1.Debug("https server : %v", addr)
	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	server.ListenAndServeTLS(crt, key)
}

func SetJsonRep(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
}

func SetHtmlRep(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
}

func SetTextRep(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
}

func FailureResp(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	logger1.Error("%v - %v", code, msg)
	io.WriteString(w, msg)
}

func SuccessRespJson(w http.ResponseWriter, r any) {
	SetJsonRep(w)
	result := model.RespSuccess(r)
	io.WriteString(w, result.String())
}

func RespResult(w http.ResponseWriter, r model.RespModel) {
	SetJsonRep(w)
	io.WriteString(w, r.String())
}
