package model

import (
	"gitee.com/vincent78/gcutil/global"
	"gitee.com/vincent78/gcutil/model"
)

var ErrorHttpParams = func(msg string) *model.ErrorModel {
	return model.NewErrModelByStr(global.ErrorCodeHttpBase+100, "the params is not i want %v", msg)
}

var ErrorHttpPanic = func(err error) *model.ErrorModel {
	return model.NewErrModel(global.ErrorCodeHttpBase+1, err)
}

var ErrorHttp404 = func(url string) *model.ErrorModel {
	return model.NewErrModelByStr(global.ErrorCodeHttpBase+2, "404 method [%v]", url)
}
