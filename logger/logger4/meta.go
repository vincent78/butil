package logger4

import "context"

type MetaExtractor func(ctx context.Context, key string) (any, bool)

var metaExtractor MetaExtractor = func(ctx context.Context, key string) (any, bool) {
	if ctx == nil {
		return nil, false
	}
	v := ctx.Value(key)
	return v, v != nil
}

func SetMetaExtractor(extractor MetaExtractor) {
	if extractor == nil {
		return
	}
	metaExtractor = extractor
}
