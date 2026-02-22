package logger4

import (
	"context"

	"github.com/vincent78/butil/config"
)

type GroupLogger struct {
	Name    string
	Loggers []*NormalLogger
}

func NewGroupLogger(name string, confs []config.LoggerConfig) *GroupLogger {
	var ls []*NormalLogger

	for _, conf := range confs {
		l := NewNormalLogger(conf)
		ls = append(ls, l)
	}

	return &GroupLogger{
		Name:    name,
		Loggers: ls,
	}
}

func (g *GroupLogger) Debug(msg string, fields ...Field) {
	for _, l := range g.Loggers {
		l.Debug(msg, fields...)
	}
}

func (g *GroupLogger) Info(msg string, fields ...Field) {
	for _, l := range g.Loggers {
		l.Info(msg, fields...)
	}
}

func (g *GroupLogger) Warn(msg string, fields ...Field) {
	for _, l := range g.Loggers {
		l.Warn(msg, fields...)
	}
}

func (g *GroupLogger) Error(msg string, fields ...Field) {
	for _, l := range g.Loggers {
		l.Error(msg, fields...)
	}
}

func (g *GroupLogger) DPanic(msg string, fields ...Field) {
	for _, l := range g.Loggers {
		l.DPanic(msg, fields...)
	}
}

func (g *GroupLogger) Panic(msg string, fields ...Field) {
	for _, l := range g.Loggers {
		l.Panic(msg, fields...)
	}
}

func (g *GroupLogger) Fatal(msg string, fields ...Field) {
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

func (l *GroupLogger) WithFields(field ...Field) {
	for _, l := range l.Loggers {
		l.WithFields(field...)
	}
}

func (l *GroupLogger) WithMap(mps map[string]any) {
	for _, l := range l.Loggers {
		l.WithMap(mps)
	}
}

func (l *GroupLogger) WithMetaCtx(ctx context.Context, keys ...string) {
	for _, l := range l.Loggers {
		l.WithMetaCtx(ctx, keys...)
	}
}
