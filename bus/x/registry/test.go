package registry

import (
	"github.com/vincent78/butil/bus/core/test"
)

type NewTest func(opts ...test.Option) test.Test

type testRegistry struct {
	Registry[NewTest]
}

func (r *testRegistry) Register(name string, v NewTest) error {
	if err := r.Registry.Register(name, v); err != nil {
		//logger.Default().Fatal(err)
	}
	return nil
}
