package gchttp

import (
	"context"
	"github.com/vincent78/butil/model"
)

type Task struct {
	Context  context.Context
	Token    string
	Method   string
	Header   map[string]string
	Url      string
	Body     string
	RespChan chan *model.RespModel
}

type TaskOption interface {
	apply(*Task)
}

type taskOptionFunc func(*Task)

func (f taskOptionFunc) apply(task *Task) {
	f(task)
}

func WithToken(token string) TaskOption {
	return taskOptionFunc(func(task *Task) {
		task.Token = token
	})
}

func WithMethod(method string) TaskOption {
	return taskOptionFunc(func(task *Task) {
		task.Method = method
	})
}

func WithBody(body string) TaskOption {
	return taskOptionFunc(func(task *Task) {
		task.Body = body
	})
}

func WithHeader(key, value string) TaskOption {
	return taskOptionFunc(func(task *Task) {
		if task.Header == nil {
			task.Header = make(map[string]string)
		}
		task.Header[key] = value
	})
}

func WithHeaderMap(m map[string]string) TaskOption {
	return taskOptionFunc(func(task *Task) {
		task.Header = m
	})
}

func NewTask(ctx context.Context, url string, respCh chan *model.RespModel, options ...TaskOption) *Task {
	if len(url) == 0 {
		panic("url cannot be empty")
	}

	task := &Task{
		Context:  ctx,
		Url:      url,
		RespChan: respCh,
		Header: map[string]string{
			"Content-Type": "application/json",
		},
	}
	for _, option := range options {
		option.apply(task)
	}
	return task
}
