package gcgin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vincent78/butil/model"
	model2 "github.com/vincent78/butil/net/gcgin/model"
)

func DoBusAction[T any](g *gin.Context, req T, action func(ctx context.Context, req T) (any, *model.ErrorModel)) {
	err := g.ShouldBindJSON(req)
	if err != nil {
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			g.JSON(http.StatusOK, model.RespWithErrModel(model2.ErrorHttpParams(errs.Error())))
		} else {
			// validator.ValidationErrors类型错误则进行翻译
			g.JSON(http.StatusOK, model.RespWithErrModel(model2.ErrorHttpParams(errs.Error())))
		}
	} else {
		data, errModel := action(g.Request.Context(), req)
		if errModel != nil {
			if data == nil {
				g.JSON(http.StatusOK, model.RespWithStr(errModel.Code, errModel.Message))
			} else {
				errModel.Data = data
				g.JSON(http.StatusOK, model.RespWithErrModelObj(errModel))
			}
		} else {
			g.JSON(http.StatusOK, model.RespSuccess(data))
		}
	}
}
