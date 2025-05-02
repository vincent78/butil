package server

import (
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/logger"
	"github.com/vincent78/butil/net/gcws/common"
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
