package gchttp

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	cl "github.com/vincent78/butil/bus/core/logger"
	logger "github.com/vincent78/butil/logger/logger4"
	"github.com/vincent78/butil/model"
	"github.com/vincent78/butil/token"
	"github.com/vincent78/butil/utils/netUtil"
	"github.com/vincent78/butil/utils/objUtil"
	"github.com/vincent78/butil/utils/strUtil"
)

var defaultTimeout = 60 * time.Second

var clientCache = make(map[string]*HttpClient)

func init() {
	clientCache = make(map[string]*HttpClient)
}

/************************************************************************
 *
 *  interface
 *
 **************************************************************************/

func GetRequest(task *Task) model.RespModel {
	task.Method = "GET"
	client := NewHttpClient("")
	return Done(task, client)
}

func HeadRequest(task *Task) model.RespModel {
	//return Request("HEAD", url, header, "")
	task.Method = "HEAD"
	client := &http.Client{}
	resp, err := client.Head(task.Url)
	if err != nil {
		return model.RespWithError(err, 500)
	} else {
		return model.RespSuccess(resp)
	}
}

func PostRequest(task *Task) model.RespModel {
	task.Method = "POST"
	if task.Header == nil {
		task.Header = make(map[string]string)
	}
	task.Header["Content-Type"] = "application/json"
	return Done(task, nil)
}

/************************************************************************
 *
 *  client
 *
 **************************************************************************/

type HttpClient struct {
	BaseUrl            string
	InsecureSkipVerify bool
	Proxy              string
	Timeout            time.Duration
	Running            bool
	client             *http.Client
	Lock               sync.RWMutex
	Logger             cl.ILoggerMethod
}

type ClientOption interface {
	apply(client *HttpClient)
}

type clientOptionFunc func(client *HttpClient)

func (f clientOptionFunc) apply(client *HttpClient) {
	f(client)
}

func WithInsecureSkipVerify(skip bool) ClientOption {
	return clientOptionFunc(func(client *HttpClient) {
		client.InsecureSkipVerify = skip
	})
}

func WithProxy(proxy string) ClientOption {
	return clientOptionFunc(func(client *HttpClient) {
		client.Proxy = proxy
	})
}
func WithTimeout(timeout time.Duration) ClientOption {
	return clientOptionFunc(func(client *HttpClient) {
		client.Timeout = timeout
	})
}
func WithLogger(log cl.ILoggerMethod) ClientOption {
	return clientOptionFunc(func(client *HttpClient) {
		client.Logger = log
	})
}

func NewHttpClient(baseUrl string, opt ...ClientOption) *HttpClient {
	client := &HttpClient{
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
			Timeout: defaultTimeout,
		},
		Timeout: defaultTimeout,
		Logger:  logger.DefaultNormalLogger(),
	}
	client.BaseUrl = netUtil.GetUrlPrefix(baseUrl)
	for _, f := range opt {
		f.apply(client)
	}

	if client.InsecureSkipVerify {
		client.client.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	if len(client.Proxy) > 0 {
		proxyURL, _ := url.Parse(client.Proxy)
		client.client.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyURL)
	}

	client.client.Timeout = client.Timeout

	return client
}

func GetClient(baseUrl string, opt ...ClientOption) *HttpClient {
	bu := netUtil.GetUrlPrefix(baseUrl)
	if c, ok := clientCache[bu]; ok {
		return c
	} else {
		c2 := NewHttpClient(baseUrl, opt...)
		clientCache[bu] = c2
		return c2
	}
}

/************************************************************************
 *
 *  core
 *
 **************************************************************************/
func DoneInChannel(task *Task, client *HttpClient) {
	go func() {
		if task.RespChan != nil {
			task.RespChan <- Done(task, client)
		} else {
			panic("the task resp channel is nil")
		}
	}()
}

func Done(task *Task, client *HttpClient) model.RespModel {
	if client == nil {
		client = NewHttpClient(task.Url)
	}

	client.Lock.Lock()
	defer client.Lock.Unlock()

	client.Running = true
	defer func() {
		client.Running = false
	}()

	tk := task.Token
	if len(tk) == 0 {
		tk = token.UniqueId()
	}

	logIn := client.Logger.WithFields(
		logger.String("token", tk),
		logger.String("method", task.Method),
		logger.String("url", task.Url),
	)

	if task.Header != nil && len(task.Header) > 0 {
		for k, v := range task.Header {
			logIn = logIn.WithFields(logger.String(k, v))
		}
	}

	if len(task.Body) > 0 {
		logIn = logIn.WithFields(logger.String("body", task.Body))
	}

	logIn.Debug("-->> http")

	var bodyReader io.Reader
	if task.Body != "" {
		bodyReader = strings.NewReader(task.Body)
	}
	logOut := client.Logger.WithFields(logger.String("token", tk))

	req, err := http.NewRequest(task.Method, task.Url, bodyReader)
	if err != nil {
		logOut.Error("<<-- http", logger.String("error", err.Error()))
		return model.RespWithError(err)
	}

	//req.Header.Set("Content-Type", "application/json")
	if task.Header != nil && len(task.Header) > 0 {
		for k, v := range task.Header {
			req.Header.Set(k, v)
		}
	}

	req = req.WithContext(task.Context)
	resp, err := client.client.Do(req)

	if err != nil {
		logOut.Error("<<-- http", logger.String("error", err.Error()))
		return model.RespWithError(err, 4000)
	}
	defer func(Body io.ReadCloser) {
		dfErr := Body.Close()
		if dfErr != nil {
			logOut.Error("<<-- http",logger.Any("statusCode", resp.StatusCode), logger.String("error", err.Error()))
		}
	}(resp.Body)
	//rep, err := ioutil.ReadAll(resp.Body)
	rep, err := io.ReadAll(resp.Body)
	if err != nil {
		logOut.Error("<<-- http", logger.Any("statusCode", resp.StatusCode),logger.String("error", err.Error()))
		return model.RespWithError(err, 4000)
	} else if strings.HasPrefix(resp.Header.Get("content-type"), "application/json") {
		r := objUtil.Parse(rep)
		logOut.Debug("<<-- http",logger.Any("statusCode", resp.StatusCode), logger.Any("response", r))
		return model.RespSuccess(r)
	} else {
		logOut.Debug("<<-- http", logger.Any("statusCode", resp.StatusCode), logger.Any("response", rep))
		return model.RespSuccess(rep)
	}
}

/************************************************************************
 *
 *  other
 *
 **************************************************************************/

func PostFormData(urlStr string, header map[string]string, body map[string]any) model.RespModel {
	if header == nil {
		header = make(map[string]string)
	}
	header["Content-Type"] = "application/x-www-form-urlencoded"

	// 用url.values方式构造form-data参数
	formValues := url.Values{}
	for k, v := range body {
		formValues.Set(k, strUtil.ToStr(v))
	}
	formDataStr := formValues.Encode()
	formDataBytes := []byte(formDataStr)
	formBytesReader := bytes.NewReader(formDataBytes)

	//生成post请求
	client := &http.Client{}
	req, err := http.NewRequest("POST", urlStr, formBytesReader)
	if err != nil {
		// handle error
		return model.RespWithError(err, 500)
	} else {
		if resp, err := client.Do(req); err != nil {
			return model.RespWithError(err, 500)
		} else {
			return model.RespSuccess(resp)
		}

	}
}
