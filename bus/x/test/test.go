package test

import (
	"context"

	"github.com/vincent78/butil/bus/core/metadata"
	ct "github.com/vincent78/butil/bus/core/test"
	"github.com/vincent78/butil/bus/x/registry"
)

func init() {
	/*
		在 Go 程序启动时，所有的 init() 函数会先于 main() 执行。
		动作：你将 NewTest 这个函数本身（作为变量或接口）存入 sync.Map。
		注意：此时 NewTest 并没有被执行，它只是像一个“函数指针”一样被存在了内存里。
	*/

	_ = registry.TestRegistry().Register("test", NewTest)
}

type test struct {
	name string
	now  string
	ops  ct.Options
}

func NewTest(opts ...ct.Option) ct.Test {
	/*
		注意： 这里只有当调用Get方法，从sync.map中load时才会调用到
	*/
	options := ct.Options{}
	for _, opt := range opts {
		opt(&options)
	}
	return &test{
		ops: options,
	}
}

func (t *test) Init(md metadata.Metadata) error {
	t.ops.Metadata.Set("xxx", md.Get("xxx"))
	return nil
}

func (t *test) Write(_ context.Context, str string, ops ...ct.Option) error {
	return nil
}
