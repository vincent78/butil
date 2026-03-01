package test

import (
	"context"

	"github.com/vincent78/butil/bus/core/metadata"
)

type Test interface {
	Init(metadata metadata.Metadata) error
	Write(context.Context, string, ...Option) error
}
