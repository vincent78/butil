package server

import (
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/logger/logger1"
	"github.com/vincent78/butil/net/gcws/common"
)

type LoginCmdInput struct {
	Name   string      `json:"name"`
	Passwd string      `json:"passwd"`
	Data   interface{} `json:"data,omitempty"`
}

type LoginCmdHandler struct {
}

func (h *LoginCmdHandler) DoAction(c *WSServer, cmd *common.WSCmd) {
	input := &LoginCmdInput{}
	em := common.ParseCmdInput(cmd.Data, input)
	if em != nil {
		logger1.ErrorByName(global.LogFileWSSName, "login cmd parse error: %v", em)
		c.Send(cmd.FailureByErrModel(em))
	} else {
		logger1.DebugByName(global.LogFileWSSName, "do login[%+v] action", input)
		c.ID = input.Name
		HubManager.Login(c)
		c.Send(cmd.Success(nil))
	}
}
