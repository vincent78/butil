package server

import (
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/model"
)

const (
	ErrorParseCmdCode int = iota + global.ErrorCodeWSServerBase + 100
	ErrorCmdInputTypeCode
	ErrorParseCmdInputCode
	ErrorNoCmdHanderCode
)

var ErrorParseCmd = func(msg string, args ...interface{}) *model.ErrorModel {
	return model.NewErrModelByStr(ErrorParseCmdCode, msg, args...)
}

var ErrorCmdInputType = func(input interface{}) *model.ErrorModel {
	return model.NewErrModelByStr(ErrorCmdInputTypeCode, "Cmd Input[%v] error", input)
}

var ErrorParseCmdInput = func(input interface{}, err error) *model.ErrorModel {
	return model.NewErrModelByStr(ErrorParseCmdInputCode, "parse Cmd Input[%v] error:%v ", input, err)
}

var ErrorNoCmdHander = func(msg string) *model.ErrorModel {
	return model.NewErrModelByStr(ErrorNoCmdHanderCode, "not find the handler: [%v]", msg)
}
