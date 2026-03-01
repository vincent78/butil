package registry

import (
	"github.com/vincent78/butil/logger/logger4"
)

type NewLogger func(opts ...logger4.Option) logger4.LoggerMethod

type loggerRegistry struct {
	registry[NewLogger]
}

func (r *loggerRegistry) Register(name string, v NewLogger) error {
	if err := r.registry.Register(name, v); err != nil {
	}
	return nil
}
