package gchttp

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	logger "github.com/vincent78/butil/logger/logger4"
	"github.com/vincent78/butil/model"
	"github.com/vincent78/butil/utils/netUtil"
)

func TestGetBaseUrl(t *testing.T) {
	urlstr := "https://api-cn.etherscan.com/api?module=stats&action=ethprice&apikey=351ZSP4VXJANR8VRAN4D9C22TSN9BQXQP9"
	t.Logf("url prefix: %s", netUtil.GetUrlPrefix(urlstr))
	urlstr = "https://api-cn.etherscan.com:909/api?module=stats&action=ethprice&apikey=351ZSP4VXJANR8VRAN4D9C22TSN9BQXQP9"
	t.Logf("url prefix: %s", netUtil.GetUrlPrefix(urlstr))
	urlstr = "https://api-cn.etherscan.com:909"
	t.Logf("url prefix: %s", netUtil.GetUrlPrefix(urlstr))
}

func TestProxyDone(t *testing.T) {
	//urlstr := "http://ipinfo.io"
	//urlstr := "http://cip.cc"
	urlstr := "https://jsonplaceholder.typicode.com/users/1"
	//proxy := "http://127.0.0.1:7897"
	proxy := "socks5://admin:U4oDtdNW@43.163.205.250:10444"
	proxyClient := NewHttpClient(urlstr,
		WithProxy(proxy),
		WithLogger(logger.DefaultNormalLogger()),
	)
	task := NewTask(context.Background(), urlstr, nil)
	proxyResp := Done(task, proxyClient)
	if proxyResp.IsSuccess() {
		t.Log(proxyResp.GetJson())
		//t.Log(proxyResp.GetJson().Get("address.street"))
	}

	client := NewHttpClient(urlstr)
	resp := Done(task, client)
	t.Log(resp.Data)
}

func TestDoneInChannel(t *testing.T) {
	urlstr := "https://red-solitary-valley.quiknode.pro/954951fe5e8f41b83d503bf29450307900fd4a6c"
	proxy := "http://127.0.0.1:7897"
	client := NewHttpClient(urlstr, WithProxy(proxy))
	respChan := make(chan model.RespModel)
	task := NewTask(context.Background(), urlstr, respChan,
		WithMethod("POST"),
		WithHeader("Content-Type", "application/json"),
		WithBody(`{"method":"eth_blockNumber","params":[],"id":1,"jsonrpc":"2.0"}`),
	)
	DoneInChannel(task, client)
	resp := <-respChan
	t.Log(resp.Data)
}

func TestPostHttpForm(t *testing.T) {
	urlApi := "http://localhost:8080/v1/discovery"
	var contentType string = "application/x-www-form-urlencoded"

	formParam := url.Values{}
	formParam.Add("id", "1010")
	formParam.Add("test", "test")
	body := strings.NewReader(formParam.Encode())
	// 或者 body = strings.NewReader("id=1010&test=test")

	resp, err := http.Post(urlApi, contentType, body)
	if err != nil {

	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		// handle error
	}
	t.Log(string(respBody))
}

func TestPostHttpJson(t *testing.T) {
	urlApi := "http://localhost:8080/v1/discovery"
	var contentType = "application/json"
	jsonParam := `{"id":"1010"}`
	resp, err := http.Post(urlApi, contentType, strings.NewReader(jsonParam))
	if err != nil {

	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		// handle error
	}
	t.Log(string(respBody))
}
