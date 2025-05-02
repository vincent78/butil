package thread

import (
	"fmt"
	"github.com/vincent78/butil/logger"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	pool, err := NewPool(20, 5)
	if err != nil {
		logger.Error("the Pool create error : %v", err.Error())
		return
	}
	var testFunc = func(arg interface{}) error {
		println("the test func: " + arg.(string))
		time.Sleep(time.Second * 3)
		return nil
	}

	for i := 0; i < 50; i++ {
		err := pool.Submit(testFunc, fmt.Sprintf("index %v", i))
		if err != nil {
			logger.Error("the task error : %v", err.Error())
			break
		}
	}
}
