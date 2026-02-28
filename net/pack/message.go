package pack

import "time"

const defaultMsgSize = 2048

type IMessage interface {
	GetCmd() uint32 //获取cmd数据
	SetCmd(uint32)  //设置cmd数据

	SetTime(int64)
	GetTime() int64

	SetData([]byte)  //设置data数据
	GetData() []byte //获取data数据

	Reset()
}

// Message 消息(包结构:cmd data)
type Message struct {
	Cmd  uint32 `json:"cmd"`  //消息:cmd
	Time int64  `json:"time"` //时间戳(MS)
	Data []byte `json:"data"` //消息:包体
}

// NewMessage 创建一个Message消息包
func NewMessage(cmd uint32, data []byte) *Message {
	return &Message{
		Cmd:  cmd,
		Time: time.Now().UTC().UnixMilli(),
		Data: data,
	}
}

// GetCmd 获取cmd数据
func (msg *Message) GetCmd() uint32 {
	return msg.Cmd
}

// SetCmd 设置cmd数据
func (msg *Message) SetCmd(cmd uint32) {
	msg.Cmd = cmd
}

// GetTime 获取cmd数据
func (msg *Message) GetTime() int64 {
	return msg.Time
}

// SetTime 设置cmd数据
func (msg *Message) SetTime(itime int64) {
	msg.Time = itime
}

// SetData 设置data数据
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}

// GetData 获取data数据
func (msg *Message) GetData() []byte {
	return msg.Data
}

func (msg *Message) Reset() {
	clear(msg.GetData())
	msg.SetCmd(0)
	msg.SetTime(0)
}
