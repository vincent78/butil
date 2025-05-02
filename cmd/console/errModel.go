package console

import (
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/model"
)

const (
	ErrorNoHandlerCode = global.ErrorCodeConsoleBase + 1 + iota
)

var ErrorNoHandler = func(msg string) *model.ErrorModel {
	return model.NewErrModelByStr(ErrorNoHandlerCode, "no handler: [%v]", msg)
}
