package registry

import (
	"github.com/vincent78/butil/bus/core/logger"
)

type loggerRegistry struct {
	registry[logger.ILoggerMethod]
}

func (r *loggerRegistry) Register(name string, v logger.ILoggerMethod) error {
	if err := r.registry.Register(name, v); err != nil {
	}
	return nil
}
