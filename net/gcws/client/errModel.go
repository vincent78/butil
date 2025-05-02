package client

import (
	"gitee.com/vincent78/gcutil/global"
	"gitee.com/vincent78/gcutil/model"
)

const (
	ErrorWSCCmdInputTypeCode int = iota + global.ErrorCodeWSClientBase + 100
)

var ErrorWSCCmdInputType = func(input interface{}) *model.ErrorModel {
	return model.NewErrModelByStr(ErrorWSCCmdInputTypeCode, "Cmd Input[%v] error", input)
}
