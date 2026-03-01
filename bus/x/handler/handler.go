package handler

import (
	"context"

	core "github.com/vincent78/butil/bus/core/handler"
	"github.com/vincent78/butil/bus/core/metadata"
	"github.com/vincent78/butil/bus/x/registry"
)

func init() {
	_ = registry.HandlerRegistry().Register("handler", NewHandler)
}

type handler struct {
	ops core.Options
}

func NewHandler(opts ...core.Option) core.Handler[metadata.Metadata] {
	/*
		注意： 这里只有当调用Get方法，从sync.map中load时才会调用到
	*/
	options := core.Options{}
	for _, opt := range opts {
		opt(&options)
	}
	return &handler{
		ops: options,
	}
}

func (h *handler) Init(metadata.Metadata) error {
	return nil
}
func (h *handler) Handler(context.Context, metadata.Metadata, ...core.Option) error {
	return nil
}
