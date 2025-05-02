package intercepter

import (
	"gitee.com/vincent78/gcutil/global"
	"gitee.com/vincent78/gcutil/logger"
	"gitee.com/vincent78/gcutil/model"
	httpModel "gitee.com/vincent78/gcutil/net/gcgin/model"
	"gitee.com/vincent78/gcutil/utils/strUtil"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							httpRequest, _ := httputil.DumpRequest(c.Request, false)
							logger.ErrorByName(global.LogFileHttpName, "----------------------")
							logger.ErrorByName(global.LogFileHttpName, "path: %v", c.Request.URL.Path)
							logger.ErrorByName(global.LogFileHttpName, "request: %v", strUtil.Bytes2String(httpRequest))
							logger.ErrorByName(global.LogFileHttpName, "error: %v", ne.Error())
							c.JSON(http.StatusOK, model.FailureRespWithErrModel(httpModel.ErrorHttpPanic(se.Err)))
							return
						}
					} else {
						c.AbortWithStatus(http.StatusInternalServerError)
					}
				} else {
					c.JSON(http.StatusOK, model.FailureRespWithErrModel(httpModel.ErrorHttpPanic(err.(error))))
				}

			}
		}()
		c.Next()
	}
}
