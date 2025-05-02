package client

import (
	"gitee.com/vincent78/gcutil/net/gcws/common"
)

type PongCMDHandler struct {
}

func (h *PongCMDHandler) doAction(c *WSClient, cmd *common.WSCmd) {
	//logger.DebugByName(global.LogFileWSCName, "%v: receive the pong: %v", c.LogPrefix(), cmd.Data)
}
