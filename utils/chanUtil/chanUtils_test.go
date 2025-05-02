package chanUtil

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestDefaultChan1Sender1Receiver(t *testing.T) {
	ch := NewDefaultChan1Sender1Receiver(4)
	ch.Send(`console.log("这是javascript吗");`)
	ch.Send("tadakuso×osukadat√")
	ch.Send(`println!("不是js 不过有可能是rust?")`)
	ch.Send("def main():\n\tif __name__=='__main__':\n\t\tprint('这边不是python啊!')")

	// 通过闭包可以在执行循环的时候利用到额外的变量
	ch.Range(func() func(interface{}) bool {
		var count int
		return func(val interface{}) bool {
			count++
			data := val.(string)
			if data == "tadakuso×osukadat√" {
				return false
			}
			fmt.Printf("No.%d, val: %s\n", count, data)
			return true
		}
	}())
	if err := ch.TrySend("114514"); err != nil {
		if errors.Is(err, ErrFullChan) {
			fmt.Println("发送失败 管道已经满了!")
		}
	}
	fmt.Printf("receive msg succeeded, msg: %s\n", ch.Receive())
	ch.Send("Nothing")
	if v, ok := ch.ReceiveWithBoolean(); ok {
		fmt.Printf("receive msg succeeded, msg: %s\n", v)
	}
	ch.Send("titotihiro no kamigakusi")

	// 阻塞操作作用于select块
	done := func() <-chan interface{} {
		done := make(chan interface{})
		go func() {
			v, ok := ch.ReceiveWithBoolean()
			done <- [2]interface{}{v, ok}
		}()
		return done
	}()
	// 因为又多了一层管道 有一定性能损失 不sleep的话还是取不到的
	time.Sleep(time.Nanosecond * 10)
	select {
	case res := <-done:
		fmt.Println(res.([2]interface{}))
	default:
		fmt.Println("取不到")
	}

	ch.StopReceive()
	if err := ch.TrySend("true is false"); true {
		fmt.Println(errors.Is(err, ErrNoLongerReceive))
	}
}
