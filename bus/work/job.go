package work

import (
	"context"
	"time"

	"github.com/vincent78/butil/bus/core/metadata"
	"github.com/vincent78/butil/model"
)

type Job struct {
	Name   string
	Ctx    context.Context
	Cancel context.CancelFunc
	Ops    *JobOptions
	Fn     JobFn
}

type JobFn func(metadata.Metadata) model.RespModel

func NewJob(ctx context.Context, name string, f JobFn, ops ...JobOption) *Job {
	sctx, cancel := context.WithCancel(ctx)
	opts := ApplyJobOptions(ops...)
	return &Job{
		Name:   name,
		Ctx:    sctx,
		Cancel: cancel,
		Ops:    opts,
		Fn:     f,
	}
}

func (j *Job) Done() model.RespModel {

	return model.SuccessResp("")
}

/*********************************************************************
 *
 *  JobOptions
 *
 *********************************************************************/

type JobOptions struct {
	Timeout time.Duration
}

type JobOption func(o *JobOptions)

func ApplyJobOptions(opts ...JobOption) *JobOptions {
	ops := &JobOptions{}
	for _, opt := range opts {
		opt(ops)
	}
	return ops
}

func WithJobTimeout(timeout time.Duration) JobOption {
	return func(o *JobOptions) {
		o.Timeout = timeout
	}
}
