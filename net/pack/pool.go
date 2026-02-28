package pack

import (
	"sync"
)

var msgUpPool = &sync.Pool{
	New: func() any {
		return msgUpAllocate()
	},
}

func msgUpAllocate() IMessage {
	return &Message{
		Cmd:  0,
		Time: 0,
		Data: make([]byte, defaultMsgSize),
	}
}

func msgUpGetFromPool() *Message {
	return msgUpPool.Get().(*Message)
}

func msgUpPutToPool(msg IMessage) {
	if msg != nil {
		msg.Reset()
		msgUpPool.Put(msg)
	}
}

// -----------------------------------------------------------------

var msgDownPool = &sync.Pool{
	New: func() any {
		return &Message{
			Cmd:  0,
			Time: 0,
			Data: nil,
		}
	},
}

func msgDownGetFromPool() IMessage {
	return msgDownPool.Get().(*Message)
}

func msgDownPutToPool(msg IMessage) {
	if msg != nil {
		msg.SetCmd(0)
		msg.SetTime(0)
		msg.SetData(nil)
		msgDownPool.Put(msg)
	}
}
