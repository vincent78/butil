package intercepter

import (
	"context"
	"gitee.com/vincent78/gcutil/global"
	"gitee.com/vincent78/gcutil/token"
	"github.com/gin-gonic/gin"
)

func TokenInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		if v, ok := ctx.Value(global.RequestTokenKey).(string); !ok || len(v) == 0 {
			ts := c.GetHeader(global.RequestTokenKey)
			if ts == "" {
				//token = token2.NewUUID().String()
				//token = strUtil.GetRandomString(9)
				ts = token.NewShort()
			}
			n := context.WithValue(ctx, global.RequestTokenKey, ts)
			c.Request = c.Request.WithContext(n)
		}
		c.Next()
	}
}
