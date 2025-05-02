package chanUtil

import (
	"errors"
)

var (
	ErrNoLongerReceive = errors.New("receiver no longer receives value from channel")
	ErrFullChan        = errors.New("channel is full")
	ErrEmptyChan       = errors.New("channel is empty")
)

type Chan1Sender1Receiver interface {
	TrySend(val interface{}) error           // 尝试发送
	Send(val interface{})                    // 接收
	StopSend()                               // 发送方停止发送
	TryReceive() (interface{}, error)        // 尝试接收
	Receive() interface{}                    // 接收 v := <-c
	ReceiveWithBoolean() (interface{}, bool) // 接收 v, ok := <-c
	StopReceive()                            // 接收方停止接收
	Range(func(val interface{}) bool)        // for-range语法 类似于sync.Map.Range
}

type defaultChan1Sender1Receiver struct {
	buffer chan interface{}
	signal chan struct{}
}

func NewDefaultChan1Sender1Receiver(size uint) Chan1Sender1Receiver {
	return &defaultChan1Sender1Receiver{
		buffer: make(chan interface{}, size),
		signal: make(chan struct{}),
	}
}

func (c *defaultChan1Sender1Receiver) TryReceive() (interface{}, error) {
	select {
	case v := <-c.buffer:
		return v, nil
	default:
		return nil, ErrEmptyChan
	}
}
func (c *defaultChan1Sender1Receiver) Receive() interface{} {
	return <-c.buffer
}
func (c *defaultChan1Sender1Receiver) ReceiveWithBoolean() (v interface{}, ok bool) {
	v, ok = <-c.buffer
	return
}
func (c *defaultChan1Sender1Receiver) StopSend() {
	close(c.buffer)
}
func (c *defaultChan1Sender1Receiver) StopReceive() {
	close(c.signal)
}
func (c *defaultChan1Sender1Receiver) Send(v interface{}) {
	c.buffer <- v
}
func (c *defaultChan1Sender1Receiver) TrySend(v interface{}) error {
	// 接收方不再接收的优先级是最高的
	// 调用者在收到ErrNoLongerReceive之后就不应该再向管道发送数据了
	select {
	case c.buffer <- v:
		// 虽然成功发了一个数据 但接收方已经不再接收倒也没啥关系
		select {
		case <-c.signal:
			return ErrNoLongerReceive
		default:
			return nil
		}
	case <-c.signal:
		return ErrNoLongerReceive
	default:
		return ErrFullChan
	}
}
func (c *defaultChan1Sender1Receiver) Range(f func(interface{}) bool) {
	for v := range c.buffer {
		if !f(v) {
			break
		}
	}
}
