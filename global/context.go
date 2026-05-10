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
	keys := GetValueAllKey(ctx)
	nCtx := context.WithValue(ctx, metaKey, &Metadata{})
	for _, key := range keys {
		va, exist := GetCtxValue(ctx, key)
		if exist {
			SetCtxValue(nCtx, key, va)
		}
	}
	return nCtx
}

func GetValueAllKey(ctx context.Context) []string {
	var keys []string
	if m, ok := ctx.Value(metaKey).(*Metadata); ok {
		m.data.Range(func(k, v interface{}) bool {
			keys = append(keys, k.(string))
			return true
		})
	}
	return keys
}

// 设置值 (即使在子协程中，由于是指针且使用了 sync.Map，也是安全的)
func SetCtxValue(ctx context.Context, key string, val any) {
	m, ok := ctx.Value(metaKey).(*Metadata)
	if !ok || m == nil {
		panic("context not include MetaData")
	} else {
		m.data.Store(key, val)
	}
}

// 获取值
func GetCtxValue(ctx context.Context, key string) (any, bool) {
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
