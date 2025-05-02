package server

import (
	"gitee.com/vincent78/gcutil/global"
	"gitee.com/vincent78/gcutil/logger"
	"gitee.com/vincent78/gcutil/net/gcws/common"
)

type LogoutCmdInput struct {
	Id string `json:"ID"`
}

type LogoutCmdHandler struct {
}

func (h *LogoutCmdHandler) DoAction(c *WSServer, _ *common.WSCmd) {
	logger.DebugByName(global.LogFileWSSName, c.NormalLogger("do logout action"))
	HubManager.Logout(c)
}
