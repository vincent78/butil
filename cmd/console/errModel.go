package console

import (
	"gitee.com/vincent78/gcutil/global"
	"gitee.com/vincent78/gcutil/model"
)

const (
	ErrorNoHandlerCode = global.ErrorCodeConsoleBase + 1 + iota
)

var ErrorNoHandler = func(msg string) *model.ErrorModel {
	return model.NewErrModelByStr(ErrorNoHandlerCode, "no handler: [%v]", msg)
}
