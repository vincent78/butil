package logger3

import (
	"fmt"
	"sync"

	"github.com/vincent78/butil/utils/strUtil"
	"github.com/vincent78/butil/utils/timeUtil"
)

type MessageFormat string

const (
	MsgFormatDefault MessageFormat = "default"
	MsgFormatJson                  = "json"
)

type MsgModel struct {
	Now     string                 `json:"ts"`
	Level   LoggerLeverl           `json:"level"`
	Format  MessageFormat          `json:"-"`
	Message string                 `json:"msg"`
	Fields  map[string]interface{} `json:"field,omitempty"`
}

var msgObjPool = &sync.Pool{New: func() any {
	return MsgModel{}
}}

func NewMsgModel(l LoggerLeverl, msg string, field map[string]interface{}) *MsgModel {
	obj := msgObjPool.Get().(MsgModel)
	obj.Now = timeUtil.NowFmtStr(timeUtil.LongFormat)
	obj.Level = l
	obj.Message = msg
	obj.Fields = field
	return &obj
}

func (msg *MsgModel) SendBack() {
	msgObjPool.Put(*msg)
}

func (msg *MsgModel) Str() string {
	m := ""
	if msg.Format == MsgFormatJson {
		m = strUtil.ToJsonStr(msg)
	} else {
		if msg.Fields != nil && len(msg.Fields) > 0 {
			m = fmt.Sprintf("%v\t %v\t %v %v", msg.Now, msg.Level.str(), msg.Message, strUtil.ToJsonStr(msg.Fields))
		} else {
			m = fmt.Sprintf("%v\t %v\t %v", msg.Now, msg.Level.str(), msg.Message)
		}
	}
	return m
}

func (msg *MsgModel) Bytes() []byte {
	str := msg.Str()
	return strUtil.ToBytes(str)
}
