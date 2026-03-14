package logger4

import (
	"context"
	"strings"

	"github.com/vincent78/butil/bus/core/logger"
	"github.com/vincent78/butil/bus/x/registry"
	"github.com/vincent78/butil/config"
)

type GroupLogger struct {
	Name    string
	Loggers []logger.ILoggerMethod
}

func NewGroupLogger(name string, confs map[string]config.LoggerConfig) *GroupLogger {
	var ls []logger.ILoggerMethod

	for n, conf := range confs {
		if strings.HasPrefix(n, name) {
			l := registry.LoggerRegistry().Get(n)
			if l == nil {
				l = NewNormalLogger(conf)
			}
			ls = append(ls, l)
		}
	}

	return &GroupLogger{
		Name:    name,
		Loggers: ls,
	}
}

func (g *GroupLogger) Debug(msg string, fields ...logger.Field) {
	for _, l := range g.Loggers {
		l.Debug(msg, fields...)
	}
}

func (g *GroupLogger) Info(msg string, fields ...logger.Field) {
	for _, l := range g.Loggers {
		l.Info(msg, fields...)
	}
}

func (g *GroupLogger) Warn(msg string, fields ...logger.Field) {
	for _, l := range g.Loggers {
		l.Warn(msg, fields...)
	}
}

func (g *GroupLogger) Error(msg string, fields ...logger.Field) {
	for _, l := range g.Loggers {
		l.Error(msg, fields...)
	}
}

func (g *GroupLogger) DPanic(msg string, fields ...logger.Field) {
	for _, l := range g.Loggers {
		l.DPanic(msg, fields...)
	}
}

func (g *GroupLogger) Panic(msg string, fields ...logger.Field) {
	for _, l := range g.Loggers {
		l.Panic(msg, fields...)
	}
}

func (g *GroupLogger) Fatal(msg string, fields ...logger.Field) {
	for _, l := range g.Loggers {
		l.Fatal(msg, fields...)
	}
}

func (g *GroupLogger) Sync() error {
	for _, l := range g.Loggers {
		err := l.Sync()
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *GroupLogger) WithFields(field ...logger.Field) logger.ILoggerMethod {
	gl := &GroupLogger{
		Name:    g.Name,
		Loggers: make([]logger.ILoggerMethod, 0, len(g.Loggers)),
	}
	for _, l := range g.Loggers {
		gl.Loggers = append(gl.Loggers, l.WithFields(field...))
	}

	return gl
}

func (g *GroupLogger) WithMap(mps map[string]any) logger.ILoggerMethod {
	gl := &GroupLogger{
		Name:    g.Name,
		Loggers: make([]logger.ILoggerMethod, 0, len(g.Loggers)),
	}
	for _, l := range g.Loggers {
		gl.Loggers = append(gl.Loggers, l.WithMap(mps))
	}

	return gl
}

func (g *GroupLogger) WithMetaCtx(ctx context.Context, keys ...string) logger.ILoggerMethod {
	gl := &GroupLogger{
		Name:    g.Name,
		Loggers: make([]logger.ILoggerMethod, 0, len(g.Loggers)),
	}
	for _, l := range g.Loggers {
		gl.Loggers = append(gl.Loggers, l.WithMetaCtx(ctx, keys...))
	}

	return gl
}

func (g *GroupLogger) WithCallerSkip(skip int) logger.ILoggerMethod {
	gl := &GroupLogger{
		Name:    g.Name,
		Loggers: make([]logger.ILoggerMethod, 0, len(g.Loggers)),
	}
	for _, l := range g.Loggers {
		gl.Loggers = append(gl.Loggers, l.WithCallerSkip(skip))
	}

	return gl
}

func (g *GroupLogger) AppendLogger(logger logger.ILoggerMethod) {
	g.Loggers = append(g.Loggers, logger)
}
