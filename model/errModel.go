package model

import (
	"fmt"
	"gitee.com/vincent78/gcutil/global"
)

type ErrorModel struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Suggest interface{} `json:"suggest,omitempty"`
	Params  interface{} `json:"params,omitempty"`
}

func NewErrModelByStr(code int, msg string, args ...interface{}) *ErrorModel {
	obj := &ErrorModel{}
	obj.Code = code
	obj.Message = fmt.Sprintf(msg, args...)
	return obj
}

func NewErrModel(code int, err error) *ErrorModel {
	str := err.Error()
	return NewErrModelByStr(code, str)
}

func (eo *ErrorModel) ToString(args ...interface{}) string {
	return fmt.Sprintf(eo.Message, args...)
}

var ErrorBaseNormal = func(info string) *ErrorModel {
	return NewErrModelByStr(global.ErrorCodeBase, "error : %v", info)
}
