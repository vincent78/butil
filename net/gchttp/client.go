package gchttp

import (
	"bytes"
	"context"
	"crypto/tls"
	logger "github.com/vincent78/butil/logger/logger4"
	"github.com/vincent78/butil/model"
	"github.com/vincent78/butil/token"
	"github.com/vincent78/butil/utils/jsonUtil"
	"github.com/vincent78/butil/utils/strUtil"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var defaultTimeout = 3 * time.Second

func GetRequest(url string, header map[string]string, timeout time.Duration) model.RespModel {
	return Request("GET", url, header, "", timeout)
}

func HeadRequest(url string) model.RespModel {
	//return Request("HEAD", url, header, "")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Head(url)
	if err != nil {
		return model.FailureRespWithError(500, err)
	} else {
		return model.SuccessResp(resp)
	}
}

func PostRequest(url string, header map[string]string, body string, timeout time.Duration) model.RespModel {
	if header == nil {
		header = make(map[string]string)
	}
	header["Content-Type"] = "application/json"
	return Request("POST", url, header, body, timeout)
}

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

func Request(method, url string, header map[string]string, body string, timeout time.Duration) model.RespModel {
	tk := ""
	//hjs, _ := json.Marshal(header)
	logInFields := make([]logger.Field, 0)

	if v, ok := header["Bus_Token"]; ok {
		tk = v
	} else {
		tk = token.UniqueId()
	}
	logInFields = append(logInFields, logger.String("token", tk))

	logInFields = append(logInFields, logger.String("method", method),
		logger.String("url", url))

	if _, ok := header["Api-Key"]; ok {
		logInFields = append(logInFields, logger.String("apiKey", header["Api-Key"]))
	}

	if len(body) > 0 {
		logInFields = append(logInFields, logger.String("body", body))
	}

	logOutFields := []logger.Field{
		logger.String("token", tk),
	}

	logger.Debug("-->> http", logInFields...)
	//跳过证书校验
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	// 使用系统的代理
	http.DefaultTransport.(*http.Transport).Proxy = http.ProxyFromEnvironment

	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		logOutFields = append(logOutFields, logger.String("error", err.Error()))
		logger.Error("<<-- http", logOutFields...)
		return model.FailureRespWithError(4000, err)
	}

	//req.Header.Set("Content-Type", "application/json")
	if header != nil && len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	//超时处理
	ctx, cancel := context.WithTimeout(req.Context(), timeout)
	defer cancel()
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logOutFields = append(logOutFields, logger.String("error", err.Error()))
		logger.Error("<<-- http", logOutFields...)
		return model.FailureRespWithError(4000, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logOutFields = append(logOutFields, logger.String("error", err.Error()))
			logger.Error("<<-- http: close the body err", logOutFields...)
		}
	}(resp.Body)
	//rep, err := ioutil.ReadAll(resp.Body)
	rep, err := io.ReadAll(resp.Body)
	if err != nil || len(rep) == 0 {
		logOutFields = append(logOutFields, logger.String("error", err.Error()))
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
