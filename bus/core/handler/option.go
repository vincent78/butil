package handler

import (
	"github.com/vincent78/butil/bus/core/metadata"
)

type Options struct {
	Metadata metadata.Metadata
}

type Option func(opts *Options)

func WithMetadata(metadata metadata.Metadata) Option {
	return func(opts *Options) {
		opts.Metadata = metadata
	}
}
