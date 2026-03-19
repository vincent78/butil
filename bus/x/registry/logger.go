package registry

import (
	"github.com/vincent78/butil/bus/core/logger"
)

type loggerRegistry struct {
	Registry[logger.ILoggerMethod]
}

func (r *loggerRegistry) Register(name string, v logger.ILoggerMethod) error {
	if err := r.Registry.Register(name, v); err != nil {
		return err
	}
	return nil
}
