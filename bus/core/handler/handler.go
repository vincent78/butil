package handler

import (
	"context"

	"github.com/vincent78/butil/bus/core/metadata"
)

type Handler[T any] interface {
	Init(metadata.Metadata) error
	Handler(context.Context, T, ...Option) error
}
