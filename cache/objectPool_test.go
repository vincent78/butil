package cache

import (
	"fmt"
	"testing"
)

func TestObjPool(t *testing.T) {
	num := func() interface{} {
		return 10.0
	}

	pool := NewObjPool(num)
	object := pool.Acquire()

	fmt.Println(pool.Inuse, "... Pool Inuse")
	fmt.Println(pool.Available, "... Pool Available")

	//assert.Equal(t, len(pool.Inuse), num())

	pool.Release(object)

	fmt.Println(pool.Inuse, "... Pool Inuse")
	fmt.Println(pool.Available, "... Pool Available")
	//assert.Equal(t, len(pool.Available), 0)
}
