package model

import (
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/model"
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
