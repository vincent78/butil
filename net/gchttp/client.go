package gchttp

import (
	"bytes"
	"crypto/tls"
	logger "github.com/vincent78/butil/logger/logger4"
	"github.com/vincent78/butil/model"
	"github.com/vincent78/butil/token"
	"github.com/vincent78/butil/utils/jsonUtil"
	"github.com/vincent78/butil/utils/netUtil"
	"github.com/vincent78/butil/utils/strUtil"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var defaultTimeout = 3 * time.Second

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
	return done(task, client)
}

func HeadRequest(task *Task) model.RespModel {
	//return Request("HEAD", url, header, "")
	task.Method = "HEAD"
	client := &http.Client{}
	resp, err := client.Head(task.Url)
	if err != nil {
		return model.FailureRespWithError(500, err)
	} else {
		return model.SuccessResp(resp)
	}
}

func PostRequest(task *Task) model.RespModel {
	task.Method = "POST"
	if task.Header == nil {
		task.Header = make(map[string]string)
	}
	task.Header["Content-Type"] = "application/json"
	return done(task, nil)
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

func NewHttpClient(baseUrl string, opt ...ClientOption) *HttpClient {
	client := &HttpClient{
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
			Timeout: defaultTimeout,
		},
		Timeout: defaultTimeout,
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
func doneInChannel(task *Task, client *HttpClient) {
	go func() {
		resp := done(task, client)
		if task.RespChan != nil {
			task.RespChan <- &resp
		} else {
			panic("the task resp channel is nil")
		}
	}()
}

func done(task *Task, client *HttpClient) model.RespModel {
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

	logInFields := make([]logger.Field, 0)
	logInFields = append(logInFields,
		logger.String("token", tk),
		logger.String("method", task.Method),
		logger.String("url", task.Url))

	if task.Header != nil && len(task.Header) > 0 {
		for k, v := range task.Header {
			logInFields = append(logInFields, logger.String(k, v))
		}
	}

	if len(task.Body) > 0 {
		logInFields = append(logInFields, logger.String("body", task.Body))
	}

	logger.Debug("-->> http", logInFields...)

	var bodyReader io.Reader
	if task.Body != "" {
		bodyReader = strings.NewReader(task.Body)
	}
	logOutFields := make([]logger.Field, 0)
	logOutFields = append(logOutFields, logger.String("token", tk))

	req, err := http.NewRequest(task.Method, task.Url, bodyReader)
	if err != nil {
		logOutFields = append(logOutFields, logger.String("error", err.Error()))
		logger.Error("<<-- http", logOutFields...)
		return model.FailureRespWithError(4000, err)
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
		logOutFields = append(logOutFields, logger.String("error", err.Error()))
		logger.Error("<<-- http", logOutFields...)
		return model.FailureRespWithError(4000, err)
	}
	defer func(Body io.ReadCloser) {
		dfErr := Body.Close()
		if dfErr != nil {
			logOutFields = append(logOutFields, logger.String("error", dfErr.Error()))
			logger.Error("<<-- http: close the body err", logOutFields...)
		}
	}(resp.Body)
	//rep, err := ioutil.ReadAll(resp.Body)
	rep, err := io.ReadAll(resp.Body)
	if err != nil {
		logOutFields = append(logOutFields, logger.Any("error", err.Error()))
		logger.Error("<<-- http", logOutFields...)
		return model.FailureRespWithError(4000, err)
	} else if strings.HasPrefix(resp.Header.Get("content-type"), "application/json") {
		r := jsonUtil.Parse(rep)
		logOutFields = append(logOutFields, logger.Any("response", r))
		logger.Debug("<<-- http", logOutFields...)
		return model.SuccessResp(r)
	} else {
		//rs := ToStr(resp)
		md := model.SuccessResp(resp)
		logger.Debug("<<-- http", logOutFields...)
		return md
	}
}

/************************************************************************
 *
 *  other
 *
 **************************************************************************/

func PostFormData(urlStr string, header map[string]string, body map[string]interface{}) model.RespModel {
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
		return model.FailureRespWithError(500, err)
	} else {
		if resp, err := client.Do(req); err != nil {
			return model.FailureRespWithError(500, err)
		} else {
			return model.SuccessResp(resp)
		}

	}
}
