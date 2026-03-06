package logger4

import (
	"context"
	"reflect"

	"github.com/vincent78/butil/bus/core/logger"
	"github.com/vincent78/butil/config"
	"github.com/vincent78/butil/global"
	"go.uber.org/zap"
)

const DefaultName = "default"

var defaultSugaredLogger *SugaredLogger

//////////////////////////////////////////////////////////////////////////////

type NormalLogger struct {
	logger *SimpleLogger
	conf   config.LoggerConfig
}

func NewNormalLogger(conf config.LoggerConfig) *NormalLogger {
	l, err := InitLoggerByConf(conf)
	if err != nil {
		panic(err)
	}
	return &NormalLogger{
		logger: l,
		conf:   conf,
	}
}
func (l *NormalLogger) GetLogger() *SimpleLogger {
	return l.logger
}

func (l *NormalLogger) Debug(msg string, fields ...logger.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *NormalLogger) Info(msg string, fields ...logger.Field) {
	l.logger.Info(msg, fields...)
}

func (l *NormalLogger) Warn(msg string, fields ...logger.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *NormalLogger) Error(msg string, fields ...logger.Field) {
	l.logger.Error(msg, fields...)
}

func (l *NormalLogger) DPanic(msg string, fields ...logger.Field) {
	l.logger.DPanic(msg, fields...)
}

func (l *NormalLogger) Panic(msg string, fields ...logger.Field) {
	l.logger.Panic(msg, fields...)
}

func (l *NormalLogger) Fatal(msg string, fields ...logger.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *NormalLogger) Sync() error {
	return l.logger.Sync()
}

func (l *NormalLogger) WithFields(field ...logger.Field) logger.ILoggerMethod {
	nl := l.logger.With(field...)
	return &NormalLogger{
		logger: nl,
		conf:   l.conf,
	}
}

func (l *NormalLogger) WithMap(mps map[string]any) logger.ILoggerMethod {
	var mp []logger.Field
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
	nl := l.logger.With(mp...)
	return &NormalLogger{
		logger: nl,
		conf:   l.conf,
	}
}

func (l *NormalLogger) WithMetaCtx(ctx context.Context, keys ...string) logger.ILoggerMethod {
	mp := make(map[string]any)
	for _, k := range keys {
		if v, e := global.GetCtxValue(ctx, k); e {
			mp[k] = v
		}
	}
	return l.WithMap(mp)

}

func (l *NormalLogger) WithCallerSkip(skip int) logger.ILoggerMethod {
	c := l.conf
	c.CallerSkip = skip
	return &NormalLogger{
		logger: l.logger.WithOptions(zap.AddCallerSkip(skip)),
		conf:   c,
	}
}
