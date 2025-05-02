package common

import (
	"gitee.com/vincent78/gcutil/global"
	"gitee.com/vincent78/gcutil/model"
)

const (
	ErrorWSMParseCmdCode int = iota + global.ErrorCodeWSCommonBase + 100
	ErrorWSMNoCmdName
)

var ErrorWSMParseCmd = func(msg string, args ...interface{}) *model.ErrorModel {
	return model.NewErrModelByStr(ErrorWSMParseCmdCode, msg, args...)
}

var ErrorWSSNoCmdName = func(id string) *model.ErrorModel {
	return model.NewErrModelByStr(ErrorWSMNoCmdName, "the cmd name is empty: %v", id)
}
