package model

import (
	"fmt"

	"github.com/vincent78/butil/global"
)

type ErrorModel struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Suggest any    `json:"suggest,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func NewErrModelByStr(code int, msg string, args ...any) *ErrorModel {
	obj := &ErrorModel{}
	obj.Code = code
	obj.Message = fmt.Sprintf(msg, args...)
	return obj
}

func NewErrModel(code int, err error) *ErrorModel {
	str := err.Error()
	return NewErrModelByStr(code, str)
}

func (eo *ErrorModel) ToString(args ...any) string {
	return fmt.Sprintf(eo.Message, args...)
}

var ErrorBaseNormal = func(info string) *ErrorModel {
	return NewErrModelByStr(global.ErrorCodeBase, "error : %v", info)
}
