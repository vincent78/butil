package logger2

import (
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc/grpclog"
)

type grpcLogger struct {
	zLog      *zap.Logger
	verbosity int
}

// ReplaceGRPCLoggerV2 replace grpc logger v2
func ReplaceGRPCLoggerV2(l *zap.Logger) {
	zLog := l.WithOptions(zap.AddCallerSkip(5)).With(zap.Bool("grpc_system", true))
	zzl := &grpcLogger{
		zLog:      zLog,
		verbosity: 0,
	}
	grpclog.SetLoggerV2(zzl)
}

func (l *grpcLogger) Info(args ...any) {
	l.zLog.Info(fmt.Sprint(args...))
}

func (l *grpcLogger) Infoln(args ...any) {
	l.zLog.Info(fmt.Sprint(args...))
}

func (l *grpcLogger) Infof(format string, args ...any) {
	l.zLog.Info(fmt.Sprintf(format, args...))
}

func (l *grpcLogger) Warning(args ...any) {
	l.zLog.Warn(fmt.Sprint(args...))
}

func (l *grpcLogger) Warningln(args ...any) {
	l.zLog.Warn(fmt.Sprint(args...))
}

func (l *grpcLogger) Warningf(format string, args ...any) {
	l.zLog.Warn(fmt.Sprintf(format, args...))
}

func (l *grpcLogger) Error(args ...any) {
	l.zLog.Error(fmt.Sprint(args...))
}

func (l *grpcLogger) Errorln(args ...any) {
	l.zLog.Error(fmt.Sprint(args...))
}

func (l *grpcLogger) Errorf(format string, args ...any) {
	l.zLog.Error(fmt.Sprintf(format, args...))
}

func (l *grpcLogger) Fatal(args ...any) {
	l.zLog.Fatal(fmt.Sprint(args...))
}

func (l *grpcLogger) Fatalln(args ...any) {
	l.zLog.Fatal(fmt.Sprint(args...))
}

func (l *grpcLogger) Fatalf(format string, args ...any) {
	l.zLog.Fatal(fmt.Sprintf(format, args...))
}

func (l *grpcLogger) V(level int) bool {
	return l.verbosity <= level
}
