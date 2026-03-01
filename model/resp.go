package model

import (
	"encoding/json"
	"fmt"

	"github.com/vincent78/butil/utils/strUtil"
)

const Success = 0

type RespModel struct {
	Code    int            `json:"code" example:"0"`
	Message string         `json:"message,omitempty" example:""`
	Header  map[string]any `json:"header,omitempty" example:""`
	Data    any            `json:"data,omitempty" example:""`
}

func NewRespModel(bytes []byte) RespModel {
	obj := RespModel{}
	err := json.Unmarshal(bytes, &obj)
	if err != nil {
		return RespWithError(500, err)
	}
	return obj
}

func SuccessResp(data any) RespModel {
	if data == nil {
		return RespModel{
			Code: Success,
		}
	} else {
		return RespModel{
			Code: Success,
			Data: data,
		}
	}
}

func RespWithStr(code int, msg string, v ...any) RespModel {
	s := msg
	if v != nil {
		s = fmt.Sprintf(msg, v)
	}

	return RespModel{
		Code:    code,
		Message: s,
	}
}

func RespWithErrModel(obj *ErrorModel, args ...any) RespModel {
	return RespModel{
		Code:    obj.Code,
		Message: obj.ToString(args...),
	}
}

func RespWithErrModelObj(obj *ErrorModel) RespModel {
	if obj == nil {
		return RespModel{Code: Success, Data: nil}
	}
	return RespModel{
		Code:    obj.Code,
		Message: obj.Message,
		Data:    obj.Data,
	}
}

func RespWithError(code int, e error) RespModel {
	if e == nil {
		return RespModel{
			Code:    code,
			Message: "",
		}
	} else {
		return RespWithStr(code, e.Error())
	}
}
func (resp RespModel) IsSuccess() bool {
	return resp.Code == Success
}
func (resp RespModel) String() string {
	return strUtil.ToStr(resp)
}

func (resp RespModel) Bytes() []byte {
	return strUtil.ToBytes(resp)
}
