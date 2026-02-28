package global

import (
	"context"
	"sync"
	"time"
)

type Metadata struct {
	data sync.Map
}

type metadataKey struct{} // 定义一个不导出的结构体，确保全局唯一
var metaKey = metadataKey{}

// 初始化并注入一个可变的 Metadata 容器
func WithMetaData(ctx context.Context) context.Context {
	return context.WithValue(ctx, metaKey, &Metadata{})
}

// 设置值 (即使在子协程中，由于是指针且使用了 sync.Map，也是安全的)
func SetCtxValue(ctx context.Context, key string, val interface{}) {
	if m, ok := ctx.Value(metaKey).(*Metadata); ok {
		m.data.Store(key, val)
	}
}

// 获取值
func GetCtxValue(ctx context.Context, key string) (interface{}, bool) {
	if m, ok := ctx.Value(metaKey).(*Metadata); ok {
		return m.data.Load(key)
	}
	return nil, false
}

/******************************************************************
 *
 * export New functions
 *
 ******************************************************************/

func NewMetaContext() context.Context {
	return WithMetaData(context.Background())
}

func NewTimeoutContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(NewMetaContext(), timeout)
}

func NewCancelContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(NewMetaContext())
}
