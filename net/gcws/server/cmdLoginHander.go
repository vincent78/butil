package server

import (
	"gitee.com/vincent78/gcutil/global"
	"gitee.com/vincent78/gcutil/logger"
	"gitee.com/vincent78/gcutil/net/gcws/common"
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
		logger.ErrorByName(global.LogFileWSSName, "login cmd parse error: %v", em)
		c.Send(cmd.FailureByErrModel(em))
	} else {
		logger.DebugByName(global.LogFileWSSName, "do login[%+v] action", input)
		c.ID = input.Name
		HubManager.Login(c)
		c.Send(cmd.Success(nil))
	}
}
