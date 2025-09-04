package zap

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	logger2 "github.com/vincent78/butil/logger/logger2/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zaplog struct {
	cfg  zap.Config
	zap  *zap.Logger
	opts logger2.Options
	sync.RWMutex
	fields map[string]interface{}
}

func (l *zaplog) Init(opts ...logger2.Option) error {
	//var err error

	for _, o := range opts {
		o(&l.opts)
	}

	zapConfig := zap.NewProductionConfig()
	if zconfig, ok := l.opts.Context.Value(configKey{}).(zap.Config); ok {
		zapConfig = zconfig
	}

	if zcconfig, ok := l.opts.Context.Value(encoderConfigKey{}).(zapcore.EncoderConfig); ok {
		zapConfig.EncoderConfig = zcconfig
	}

	writer, ok := l.opts.Context.Value(writerKey{}).(io.Writer)
	if !ok {
		writer = os.Stdout
	}

	skip, ok := l.opts.Context.Value(callerSkipKey{}).(int)
	if !ok || skip < 1 {
		skip = 1
	}

	// Set log Level if not default
	zapConfig.Level = zap.NewAtomicLevel()
	if l.opts.Level != logger2.InfoLevel {
		zapConfig.Level.SetLevel(loggerToZapLevel(l.opts.Level))
	}
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapConfig.EncoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(writer)),
		zapConfig.Level)

	log := zap.New(logCore, zap.AddCaller(), zap.AddCallerSkip(skip), zap.AddStacktrace(zap.DPanicLevel))
	//log, err := zapConfig.Build(zap.AddCallerSkip(skip))
	//if err != nil {
	//	return err
	//}

	// Adding seed fields if exist
	if l.opts.Fields != nil {
		data := []zap.Field{}
		for k, v := range l.opts.Fields {
			data = append(data, zap.Any(k, v))
		}
		log = log.With(data...)
	}

	// Adding namespace
	if namespace, ok := l.opts.Context.Value(namespaceKey{}).(string); ok {
		log = log.With(zap.Namespace(namespace))
	}

	// defer log.Sync() ??

	l.cfg = zapConfig
	l.zap = log
	l.fields = make(map[string]interface{})

	return nil
}

func (l *zaplog) Fields(fields map[string]interface{}) logger2.Logger {
	l.Lock()
	nfields := make(map[string]interface{}, len(l.fields))
	for k, v := range l.fields {
		nfields[k] = v
	}
	l.Unlock()
	for k, v := range fields {
		nfields[k] = v
	}

	data := make([]zap.Field, 0, len(nfields))
	for k, v := range fields {
		data = append(data, zap.Any(k, v))
	}

	zl := &zaplog{
		cfg:    l.cfg,
		zap:    l.zap,
		opts:   l.opts,
		fields: nfields,
	}

	return zl
}

func (l *zaplog) Error(err error) logger2.Logger {
	return l.Fields(map[string]interface{}{"error": err})
}

func (l *zaplog) Log(level logger2.Level, args ...interface{}) {
	l.RLock()
	data := make([]zap.Field, 0, len(l.fields))
	for k, v := range l.fields {
		data = append(data, zap.Any(k, v))
	}
	l.RUnlock()

	lvl := loggerToZapLevel(level)
	msg := fmt.Sprint(args...)
	switch lvl {
	case zap.DebugLevel:
		l.zap.Debug(msg, data...)
	case zap.InfoLevel:
		l.zap.Info(msg, data...)
	case zap.WarnLevel:
		l.zap.Warn(msg, data...)
	case zap.ErrorLevel:
		l.zap.Error(msg, data...)
	case zap.FatalLevel:
		l.zap.Fatal(msg, data...)
	}
}

func (l *zaplog) Logf(level logger2.Level, format string, args ...interface{}) {
	l.RLock()
	data := make([]zap.Field, 0, len(l.fields))
	for k, v := range l.fields {
		data = append(data, zap.Any(k, v))
	}
	l.RUnlock()

	lvl := loggerToZapLevel(level)
	msg := fmt.Sprintf(format, args...)
	switch lvl {
	case zap.DebugLevel:
		l.zap.Debug(msg, data...)
	case zap.InfoLevel:
		l.zap.Info(msg, data...)
	case zap.WarnLevel:
		l.zap.Warn(msg, data...)
	case zap.ErrorLevel:
		l.zap.Error(msg, data...)
	case zap.FatalLevel:
		l.zap.Fatal(msg, data...)
	}
}

func (l *zaplog) String() string {
	return "zap"
}

func (l *zaplog) Options() logger2.Options {
	return l.opts
}

// New builds a new logger based on options
func NewLogger(opts ...logger2.Option) (logger2.Logger, error) {
	// Default options
	options := logger2.Options{
		Level:   logger2.InfoLevel,
		Fields:  make(map[string]interface{}),
		Out:     os.Stderr,
		Context: context.Background(),
	}

	l := &zaplog{opts: options}
	if err := l.Init(opts...); err != nil {
		return nil, err
	}

	return l, nil
}

func loggerToZapLevel(level logger2.Level) zapcore.Level {
	switch level {
	case logger2.TraceLevel, logger2.DebugLevel:
		return zap.DebugLevel
	case logger2.InfoLevel:
		return zap.InfoLevel
	case logger2.WarnLevel:
		return zap.WarnLevel
	case logger2.ErrorLevel:
		return zap.ErrorLevel
	case logger2.FatalLevel:
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func zapToLoggerLevel(level zapcore.Level) logger2.Level {
	switch level {
	case zap.DebugLevel:
		return logger2.DebugLevel
	case zap.InfoLevel:
		return logger2.InfoLevel
	case zap.WarnLevel:
		return logger2.WarnLevel
	case zap.ErrorLevel:
		return logger2.ErrorLevel
	case zap.FatalLevel:
		return logger2.FatalLevel
	default:
		return logger2.InfoLevel
	}
}
