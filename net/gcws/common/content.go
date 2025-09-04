package common

import (
	"time"

	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/logger/logger1"
	"github.com/vincent78/butil/timewheel"
)

// service receive cmd
const ()

// client received cmd
const ()

// client & server received cmd
const (
	CmdNameClose     = "close"
	CmdNameLogin     = "login"
	CmdNameLogout    = "logout"
	CmdNamePing      = "ping"
	CmdNamePong      = "pong"
	CmdNameConnected = "connected"
)

const (
	CmdCategorySys          = "sys"
	CmdCategoryDefault      = ""
	CmdCategoryLogicGame    = "logicGame"
	CmdCategoryLogicServer  = "logicServer"
	CmdCategoryLogicDesktop = "logicDesktop"
	CmdCategoryLogicPlayer  = "logicPlayer"
)

const (
	CmdHeaderKeyMustResp = "mustResp"
)

type MsgType rune

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

const (
	MsgTypeUnknowChar MsgType = '0'
	MsgTypeCmdChar            = '1' // msgType is cmd
	MsgTypeRespChar           = '2' // msgType is cmdResp
)

func InitTimewheel(logName string) {
	if global.Timewheel == nil {
		logger1.DebugByName(logName, "init timewheel")
		global.Timewheel, _ = timewheel.NewTimeWheel(1*time.Second, 360, timewheel.TickSafeMode())
		global.Timewheel.Start()
	}
}

func DetachTypeAndContent(source []byte) (MsgType, []byte) {
	if len(source) > 0 {
		tp := GetMsgType(source[0])
		if MsgTypeUnknowChar == tp {
			return tp, source
		} else {
			return tp, source[1:]
		}
	} else {
		return 0, nil
	}
}

func GetMsgType(source byte) MsgType {
	c := rune(source)
	if c == MsgTypeCmdChar {
		return MsgTypeCmdChar
	} else if c == MsgTypeRespChar {
		return MsgTypeRespChar
	} else {
		return MsgTypeUnknowChar
	}
}
