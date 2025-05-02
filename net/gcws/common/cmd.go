package common

import (
	"encoding/json"
	"fmt"
	"github.com/vincent78/butil/model"
	"github.com/vincent78/butil/token"
	"github.com/vincent78/butil/utils/strUtil"
)

type WSCmd struct {
	ID        string                 `json:"id"`                  // 唯一标识
	Category  string                 `json:"category,omitempty"`  // 类型
	Cmd       string                 `json:"cmd"`                 // 命令
	Data      interface{}            `json:"data,omitempty"`      // 命令的内容
	Trans     map[string]interface{} `json:"trans,omitempty"`     // 透传内容
	PID       string                 `json:"pid"`                 // 非空时，大部分为回复(resp)
	Header    map[string]string      `json:"header,omitempty"`    // 命令的头信息
	Timestamp int64                  `json:"timestamp,omitempty"` // 时间戳 （单位：毫秒）
}

type WSCmdOption struct {
	Key   string
	Value interface{}
}

func NewWSCmdOptionMustResp() WSCmdOption {
	return WSCmdOption{
		Key:   CmdHeaderKeyMustResp,
		Value: true,
	}
}

func ParseCmdInput[T any](source interface{}, obj T) *model.ErrorModel {
	if source == nil {
		return nil
	}
	bytes, err := json.Marshal(source)
	if err != nil {
		return model.ErrorBaseNormal(err.Error())
	}
	err = json.Unmarshal(bytes, obj)
	if err != nil {
		return model.ErrorBaseNormal(err.Error())
	}
	return nil
}

func NewWSCmd(cmd string, ops ...WSCmdOption) *WSCmd {
	ws := &WSCmd{
		ID:        token.NewShort(),
		Category:  CmdCategoryDefault,
		Cmd:       cmd,
		Data:      nil,
		Trans:     nil,
		Timestamp: 0,
	}

	if ops != nil && len(ops) > 0 {
		for i := 0; i < len(ops); i++ {
			op := ops[i]
			if op.Key == CmdHeaderKeyMustResp {
				ws.Trans = map[string]interface{}{CmdHeaderKeyMustResp: op.Value}
			}
		}
	}
	return ws
}

func NewWSCmdByCategory(category, cmd string, ops ...WSCmdOption) *WSCmd {
	wsCmd := NewWSCmd(cmd, ops...)
	wsCmd.Category = category
	return wsCmd
}

func ParseCmd(bytes []byte) (*WSCmd, *model.ErrorModel) {
	cmd := &WSCmd{}
	err := json.Unmarshal(bytes, &cmd)
	if err != nil {
		return cmd, ErrorWSMParseCmd(err.Error())
	} else {
		if cmd.Cmd == "" && cmd.PID == "" {
			return cmd, ErrorWSSNoCmdName(cmd.ID)
		} else {
			if cmd.ID == "" {
				cmd.ID = token.NewShort()
			}
			return cmd, nil
		}
	}
}

func (cmd *WSCmd) Bytes() []byte {
	return strUtil.ToBytes(cmd)
}

func (cmd *WSCmd) String() string {
	return strUtil.ToStr(cmd)
}

func (cmd *WSCmd) SuccessWithName(name string, data interface{}) *WSCmd {
	obj := cmd.Success(data)
	obj.Cmd = name
	return obj
}

func (cmd *WSCmd) Success(data interface{}) *WSCmd {
	return &WSCmd{
		ID:  token.NewShort(),
		Cmd: "",
		PID: cmd.ID,
		Data: model.RespModel{
			Code: model.Success,
			Data: data,
		},
	}
}

func (cmd *WSCmd) FailureWithName(name string, code int, msg string) *WSCmd {
	obj := cmd.Failure(code, msg)
	obj.Cmd = name
	return obj
}

func (cmd *WSCmd) Failure(code int, msg string) *WSCmd {
	return &WSCmd{
		ID:  token.NewShort(),
		Cmd: "",
		PID: cmd.ID,
		Data: model.RespModel{
			Code:    code,
			Message: msg,
		},
	}
}

func (cmd *WSCmd) FailureByErrModelWithName(name string, err *model.ErrorModel) *WSCmd {
	obj := cmd.FailureByErrModel(err)
	obj.Cmd = name
	return obj
}

func (cmd *WSCmd) FailureByErrModel(err *model.ErrorModel) *WSCmd {
	return &WSCmd{
		ID:  token.NewShort(),
		Cmd: "",
		PID: cmd.ID,
		Data: model.RespModel{
			Code:    err.Code,
			Message: err.Message,
		},
	}
}

func (cmd *WSCmd) HandlerKey() string {
	return fmt.Sprintf("%v:%v", cmd.Category, cmd.Cmd)
}
