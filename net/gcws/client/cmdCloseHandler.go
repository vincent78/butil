package client

import (
	"gitee.com/vincent78/gcutil/global"
	"gitee.com/vincent78/gcutil/logger"
	"gitee.com/vincent78/gcutil/net/gcws/common"
)

type CloseCmdInput struct {
}

type CloseCMDHandler struct {
}

func (h *CloseCMDHandler) doAction(c *WSClient, cmd *common.WSCmd) {
	logger.InfoByName(global.LogFileWSCName, c.NormalLogger("do close action"))
	c.Resp(cmd.Success("ok"))
	c.url = nil
	c.DisConn()
}
