package registry

import (
	"github.com/vincent78/butil/bus/core/handler"
	"github.com/vincent78/butil/bus/core/metadata"
)

type NewHandler func(opts ...handler.Option) handler.Handler[metadata.Metadata]

type handlerRegistry struct {
	Registry[NewHandler]
}

func (r *handlerRegistry) Register(name string, v NewHandler) error {
	if err := r.Registry.Register(name, v); err != nil {
		return err
	}
	return nil
}
