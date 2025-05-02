package model

import (
	"encoding/json"
	"fmt"

	"github.com/vincent78/butil/utils/strUtil"
)

type ResultModel[T any] struct {
	Code    int                    `json:"code" example:"0"`
	Message string                 `json:"message,omitempty" example:""`
	Header  map[string]interface{} `json:"header,omitempty" example:""`
	Data    T                      `json:"data,omitempty" example:""`
}

func NewResultModel[T any](bytes []byte) ResultModel[T] {
	obj := ResultModel[T]{}
	err := json.Unmarshal(bytes, &obj)
	if err != nil {
		return FailureResultWithError[T](500, err)
	}
	return obj
}

func SuccessResult[T any](data T) ResultModel[T] {
	return ResultModel[T]{
		Code: Success,
		Data: data,
	}
}

func FailureResultWithStr[T any](code int, msg string, v ...interface{}) ResultModel[T] {
	s := msg
	if v != nil {
		s = fmt.Sprintf(msg, v)
	}

	return ResultModel[T]{
		Code:    code,
		Message: s,
	}
}

func FailureResultWithErrModel[T any](obj *ErrorModel, args ...interface{}) ResultModel[T] {
	return ResultModel[T]{
		Code:    obj.Code,
		Message: obj.ToString(args...),
	}
}

func FailureResultWithErrModelObj[T any](data T, obj *ErrorModel) ResultModel[T] {

	return ResultModel[T]{
		Code:    obj.Code,
		Message: obj.Message,
		Data:    data,
	}
}

func FailureResultWithError[T any](code int, e error) ResultModel[T] {
	if e == nil {
		return ResultModel[T]{
			Code:    code,
			Message: "",
		}
	} else {
		return FailureResultWithStr[T](code, e.Error())
	}
}
func (resp ResultModel[T]) IsSuccess() bool {
	return resp.Code == Success
}
func (resp ResultModel[T]) String() string {
	return strUtil.ToStr(resp)
}

func (resp ResultModel[T]) Bytes() []byte {
	return strUtil.ToBytes(resp)
}
