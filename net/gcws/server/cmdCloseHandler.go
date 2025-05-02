package server

import (
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/logger"
	"github.com/vincent78/butil/net/gcws/common"
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
