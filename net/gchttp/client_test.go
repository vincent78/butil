package gchttp

import (
	"github.com/vincent78/butil/logger"
	"github.com/vincent78/butil/model"
	"testing"
)

func TestGetHttp(t *testing.T) {
	r := GetRequest("https://api-cn.etherscan.com/api?module=stats&action=ethprice&apikey=351ZSP4VXJANR8VRAN4D9C22TSN9BQXQP9")

	if r.Code == model.Success {
		logger.Debug("http get result. %+v ", r.Data)
	} else {
		logger.Error("http get error: %+v", r.Message)
	}

}

func TestPostHttp(t *testing.T) {
	r := PostRequest("http://10.4.62.243:6666", nil, `{"jsonrpc":"2.0","method":"web3_clientVersion","params":[],"id":67}`)
	logger.Debug("http post result. %v ", r.Message)
}

func TestHeadRequest(t *testing.T) {
	r := HeadRequest("https://api-cn.etherscan.com/api?module=stats&action=ethprice&apikey=351ZSP4VXJANR8VRAN4D9C22TSN9BQXQP9")
	if r.IsSuccess() {
		logger.Info("%v", r.String())
	} else {
		logger.Error("%v", r.String())
	}
}
