package server

import (
	"gitee.com/vincent78/gcutil/global"
	"gitee.com/vincent78/gcutil/logger"
	"gitee.com/vincent78/gcutil/net/gcws/common"
)

// 客户端主动要求断开

type CloseCmdInput struct {
}

type CloseCmdHandler struct {
}

func (h *CloseCmdHandler) DoAction(c *WSServer, _ *common.WSCmd) {
	logger.InfoByName(global.LogFileWSSName, c.NormalLogger("do close action by client"))
	HubManager.unregister <- c
}
