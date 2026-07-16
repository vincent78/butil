package logger4

import (
	"io"
	"sync"

	"github.com/vincent78/butil/bus/core/logger"
)

var loggerRegistry sync.Map

func IsRegistered(name string) bool {
	_, ok := loggerRegistry.Load(name)
	return ok
}

func unregisterLogger(name string) {
	if v, ok := loggerRegistry.Load(name); ok {
		if closer, ok := v.(io.Closer); ok {
			_ = closer.Close()
		}
		loggerRegistry.Delete(name)
	}
}

func loadLogger(name string) logger.ILoggerMethod {
	if name == "" {
		return nil
	}
	v, ok := loggerRegistry.Load(name)
	if !ok {
		return nil
	}
	l, _ := v.(logger.ILoggerMethod)
	return l
}

func storeLogger(name string, l logger.ILoggerMethod) {
	if name == "" || l == nil {
		return
	}
	loggerRegistry.Store(name, l)
}
