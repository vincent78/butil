package model

import (
	"encoding/json"
	"fmt"

	"gitee.com/vincent78/gcutil/utils/strUtil"
)

const Success = 0

type RespModel struct {
	Code    int                    `json:"code" example:"0"`
	Message string                 `json:"message,omitempty" example:""`
	Header  map[string]interface{} `json:"header,omitempty" example:""`
	Data    interface{}            `json:"data,omitempty" example:""`
}

func NewRespModel(bytes []byte) RespModel {
	obj := RespModel{}
	err := json.Unmarshal(bytes, &obj)
	if err != nil {
		return FailureRespWithError(500, err)
	}
	return obj
}

func SuccessResp(data interface{}) RespModel {
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
	//else if len(data) == 1 {
	//	return RespModel{
	//		Code: Success,
	//		Data: data[0],
	//	}
	//} else {
	//	return RespModel{
	//		Code: Success,
	//		Data: strUtil.FmtSlice("%v", " ", data...),
	//	}
	//}
}

func FailureRespWithStr(code int, msg string, v ...interface{}) RespModel {
	s := msg
	if v != nil {
		s = fmt.Sprintf(msg, v)
	}

	return RespModel{
		Code:    code,
		Message: s,
	}
}

func FailureRespWithErrModel(obj *ErrorModel, args ...interface{}) RespModel {
	return RespModel{
		Code:    obj.Code,
		Message: obj.ToString(args...),
	}
}

func FailureRespWithErrModelObj(data interface{}, obj *ErrorModel) RespModel {
	if obj == nil {
		return RespModel{Code: Success, Data: data}
	}
	return RespModel{
		Code:    obj.Code,
		Message: obj.Message,
		Data:    data,
	}
}

func FailureRespWithError(code int, e error) RespModel {
	if e == nil {
		return RespModel{
			Code:    code,
			Message: "",
		}
	} else {
		return FailureRespWithStr(code, e.Error())
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
