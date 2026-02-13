package thread

import (
	"bytes"
	"context"
	"errors"
	"runtime"
	"strconv"
	"time"
)

// 获取GOROUTINE的ID
func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// 并发协程数量
var goChanNumberTotal = make(chan struct{}, 10000)

// GO goroutine
func GO(handler func()) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	select {
	case goChanNumberTotal <- struct{}{}:
		// ok
	case _ = <-ctx.Done():
		return errors.New("too many goroutines time out")
	}

	go func() {
		defer func() {
			<-goChanNumberTotal

			if err := recover(); err != nil {
				//logger.Errorf(ctx, "GO panic, running:%d/%d, error:%v, stack:%s", len(goChanNumberTotal), cap(goChanNumberTotal), err, BytesToString(debug.Stack()))
			}
		}()

		handler()
	}()

	return nil
}
