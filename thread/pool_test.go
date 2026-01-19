package thread

import (
	"fmt"
	"testing"
	"time"

	"github.com/vincent78/butil/logger/logger1"
)

func TestPool(t *testing.T) {
	pool, err := NewPool(20, 5)
	if err != nil {
		logger1.Error("the Pool create error : %v", err.Error())
		return
	}
	var testFunc = func(arg any) error {
		println("the test func: " + arg.(string))
		time.Sleep(time.Second * 3)
		return nil
	}

	for i := 0; i < 50; i++ {
		err := pool.Submit(testFunc, fmt.Sprintf("index %v", i))
		if err != nil {
			logger1.Error("the task error : %v", err.Error())
			break
		}
	}
}
