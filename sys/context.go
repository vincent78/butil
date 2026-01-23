package sys

import (
	"context"
	"fmt"
)

func ctxString(ctx context.Context, key any) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(key); v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}
