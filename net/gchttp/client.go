package gchttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/vincent78/butil/logger/logger1"
	"github.com/vincent78/butil/model"
	"github.com/vincent78/butil/token"
	"github.com/vincent78/butil/utils/strUtil"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var timeout = 3 * time.Second

func GetRequest(url string) model.RespModel {
	return Request("GET", url, nil, "")
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

func PostRequest(url string, header map[string]string, body string) model.RespModel {
	if header == nil {
		header = make(map[string]string)
	}
	header["Content-Type"] = "application/json"
	return Request("POST", url, header, body)
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

func Request(method, url string, header map[string]string, body string) model.RespModel {
	tk := token.UniqueId()
	logger1.Debug("-->> http[%v] method:%v url:%v header:%v body:%v", tk, method, url, header, body)
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
		logger1.Debug("<<-- http[%v] %v", tk, err.Error())
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
		logger1.Debug("<<-- http[%v] %v", tk, err.Error())
		return model.FailureRespWithError(4000, err)
	}
	defer resp.Body.Close()
	//rep, err := ioutil.ReadAll(resp.Body)
	rep, err := io.ReadAll(resp.Body)
	if err != nil {
		logger1.Debug("<<-- http[%v] %v", tk, err.Error())
		return model.FailureRespWithError(4000, err)
	} else if len(rep) == 0 {
		logger1.Debug("<<-- http[%v] the response is null", tk)
		return model.FailureRespWithError(4000, err)
	} else if strings.HasPrefix(resp.Header.Get("content-type"), "application/json") {
		mapv := make(map[string]interface{})
		err = json.Unmarshal(rep, &mapv)
		if err != nil {
			logger1.Debug("<<-- http[%v] the reponse is not json!", tk)
			return model.FailureRespWithError(4000, err)
		} else {
			md := model.SuccessResp(mapv)
			logger1.Debug("<<-- http[%v] %v", tk, md)
			return md
		}
	} else {
		md := model.SuccessResp(strUtil.ToStr(rep))
		logger1.Debug("<<-- http[%v] %v", tk, md.String())
		return md
	}
}
