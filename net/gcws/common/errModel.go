package common

import (
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/model"
)

const (
	ErrorWSMParseCmdCode int = iota + global.ErrorCodeWSCommonBase + 100
	ErrorWSMNoCmdName
)

var ErrorWSMParseCmd = func(msg string, args ...any) *model.ErrorModel {
	return model.NewErrModelByStr(ErrorWSMParseCmdCode, msg, args...)
}

var ErrorWSSNoCmdName = func(id string) *model.ErrorModel {
	return model.NewErrModelByStr(ErrorWSMNoCmdName, "the cmd name is empty: %v", id)
}
