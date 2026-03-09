package model

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/vincent78/butil/utils/strUtil"
)

/**********************

func (bp *BasePlugin) ContractBalance(contract, account string) gm.ResultModel[big.Int] {
	return gm.FailureResultWithError[big.Int](bp.notImplementError(),500)
}

 ***********************/

type ResultModel[T any] struct {
	Code    int            `json:"code" example:"0"`
	Message string         `json:"message,omitempty" example:""`
	Header  map[string]any `json:"header,omitempty" example:""`
	Data    T              `json:"data,omitempty" example:""`
}

func NewResultModel[T any](bytes []byte) ResultModel[T] {
	obj := ResultModel[T]{}
	err := json.Unmarshal(bytes, &obj)
	if err != nil {
		return ResultWithError[T](err, 500)
	}
	return obj
}

func ResultSuccess[T any](data T) ResultModel[T] {
	return ResultModel[T]{
		Code: SuccessCode,
		Data: data,
	}
}

func ResultWithStr[T any](code int, msg string, v ...any) ResultModel[T] {
	s := msg
	if v != nil {
		s = fmt.Sprintf(msg, v)
	}

	return ResultModel[T]{
		Code:    code,
		Message: s,
	}
}

func ResultWithErrModel[T any](obj *ErrorModel, args ...any) ResultModel[T] {
	return ResultModel[T]{
		Code:    obj.Code,
		Message: obj.ToString(args...),
	}
}

func ResultWithErrModelObj[T any](data T, obj *ErrorModel) ResultModel[T] {

	return ResultModel[T]{
		Code:    obj.Code,
		Message: obj.Message,
		Data:    data,
	}
}

func ResultWithError[T any](e error, codes ...int) ResultModel[T] {
	code := FailureCode
	if len(codes) > 0 {
		code = codes[0]
	}
	if e == nil {
		return ResultModel[T]{
			Code:    code,
			Message: "",
		}
	} else {
		return ResultWithStr[T](code, e.Error())
	}
}
func (r ResultModel[T]) IsSuccess() bool {
	return r.Code == SuccessCode
}
func (r ResultModel[T]) String() string {
	return strUtil.ToStr(r)
}

func (r ResultModel[T]) Bytes() []byte {
	return strUtil.ToBytes(r)
}

func (r ResultModel[T]) Error() error {
	if r.IsSuccess() {
		return nil
	}
	return errors.New(r.Message)
}
