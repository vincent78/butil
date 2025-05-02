package gcgin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vincent78/butil/model"
	model2 "github.com/vincent78/butil/net/gcgin/model"
	"net/http"
)

func DoBusAction[T any](g *gin.Context, req T, action func(ctx context.Context, req T) (interface{}, *model.ErrorModel)) {
	err := g.ShouldBindJSON(req)
	if err != nil {
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			g.JSON(http.StatusOK, model.FailureRespWithErrModel(model2.ErrorHttpParams(errs.Error())))
		} else {
			// validator.ValidationErrors类型错误则进行翻译
			g.JSON(http.StatusOK, model.FailureRespWithErrModel(model2.ErrorHttpParams(errs.Error())))
		}
	} else {
		data, errModel := action(g.Request.Context(), req)
		if errModel != nil {
			if data == nil {
				g.JSON(http.StatusOK, model.FailureRespWithStr(errModel.Code, errModel.Message))
			} else {
				g.JSON(http.StatusOK, model.FailureRespWithErrModelObj(data, errModel))
			}
		} else {
			g.JSON(http.StatusOK, model.SuccessResp(data))
		}
	}
}
