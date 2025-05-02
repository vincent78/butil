package client

import (
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/model"
)

const (
	ErrorWSCCmdInputTypeCode int = iota + global.ErrorCodeWSClientBase + 100
)

var ErrorWSCCmdInputType = func(input interface{}) *model.ErrorModel {
	return model.NewErrModelByStr(ErrorWSCCmdInputTypeCode, "Cmd Input[%v] error", input)
}
