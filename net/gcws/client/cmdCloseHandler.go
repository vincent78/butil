package client

import (
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/logger/logger1"
	"github.com/vincent78/butil/net/gcws/common"
)

type CloseCmdInput struct {
}

type CloseCMDHandler struct {
}

func (h *CloseCMDHandler) doAction(c *WSClient, cmd *common.WSCmd) {
	logger1.InfoByName(global.LogFileWSCName, c.NormalLogger("do close action"))
	c.Resp(cmd.Success("ok"))
	c.url = nil
	c.DisConn()
}
