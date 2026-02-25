package logger4

import (
	"context"
	"reflect"

	"github.com/vincent78/butil/config"
	"github.com/vincent78/butil/global"
	"go.uber.org/zap"
)

type NormalLogger struct {
	logger *SimpleLogger
	conf   config.LoggerConfig
}

func NewNormalLogger(conf config.LoggerConfig) *NormalLogger {
	logger, err := InitLoggerByConf(conf)
	if err != nil {
		panic(err)
	}
	return &NormalLogger{
		logger: logger,
		conf:   conf,
	}
}

func (l *NormalLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fields...)
}

func (l *NormalLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}

func (l *NormalLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fields...)
}

func (l *NormalLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}

func (l *NormalLogger) DPanic(msg string, fields ...Field) {
	l.logger.DPanic(msg, fields...)
}

func (l *NormalLogger) Panic(msg string, fields ...Field) {
	l.logger.Panic(msg, fields...)
}

func (l *NormalLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *NormalLogger) Sync() error {
	return l.logger.Sync()
}

func (l *NormalLogger) WithFields(field ...Field) {
	l.logger = l.logger.With(field...)
}

func (l *NormalLogger) WithMap(mps map[string]any) {
	var mp []Field
	for k, v := range mps {
		if v == nil {
			continue
		}
		ty := reflect.TypeOf(v)
		if ty.Kind() == reflect.Func {
			continue
		}

		if ty.Kind() == reflect.Ptr {
			ty = ty.Elem()
		}

		if ty.Kind() == reflect.Int {
			mp = append(mp, zap.Int(k, v.(int)))
		} else if ty.Kind() == reflect.Int8 {
			mp = append(mp, zap.Int8(k, v.(int8)))
		} else if ty.Kind() == reflect.Int16 {
			mp = append(mp, zap.Int16(k, v.(int16)))
		} else if ty.Kind() == reflect.Int32 {
			mp = append(mp, zap.Int32(k, v.(int32)))
		} else if ty.Kind() == reflect.Int64 {
			mp = append(mp, zap.Int64(k, v.(int64)))
		} else if ty.Kind() == reflect.Uint {
			mp = append(mp, zap.Uint(k, v.(uint)))
		} else if ty.Kind() == reflect.Uint8 {
			mp = append(mp, zap.Uint8(k, v.(uint8)))
		} else if ty.Kind() == reflect.Uint16 {
			mp = append(mp, zap.Uint16(k, v.(uint16)))
		} else if ty.Kind() == reflect.Uint32 {
			mp = append(mp, zap.Uint32(k, v.(uint32)))
		} else if ty.Kind() == reflect.Uint64 {
			mp = append(mp, zap.Uint64(k, v.(uint64)))
		} else if ty.Kind() == reflect.Float32 {
			mp = append(mp, zap.Float32(k, v.(float32)))
		} else if ty.Kind() == reflect.Float64 {
			mp = append(mp, zap.Float64(k, v.(float64)))
		} else if ty.Kind() == reflect.Bool {
			mp = append(mp, zap.Bool(k, v.(bool)))
		} else if ty.Kind() == reflect.String {
			mp = append(mp, zap.String(k, v.(string)))
		} else if ty.Kind() == reflect.Slice && ty.Elem().Kind() == reflect.Uint8 {
			mp = append(mp, zap.ByteString(k, v.([]byte)))
		} else {
			mp = append(mp, zap.Any(k, v))
		}
	}
	l.logger = l.logger.With(mp...)
}

func (l *NormalLogger) WithMetaCtx(ctx context.Context, keys ...string) {
	mp := make(map[string]any)
	for _, k := range keys {
		if v, e := global.GetCtxValue(ctx, k); e {
			mp[k] = v
		}
	}
	if len(mp) > 0 {
		l.WithMap(mp)
	}
}

func (l *NormalLogger) WithCallerSkip(skip int) *NormalLogger {
	return &NormalLogger{
		logger: l.logger.WithOptions(zap.AddCallerSkip(skip)),
		conf:   l.conf,
	}
}
