package intercepter

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/token"
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
