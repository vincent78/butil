package client

import (
	"github.com/vincent78/butil/net/gcws/common"
)

type PongCMDHandler struct {
}

func (h *PongCMDHandler) doAction(c *WSClient, cmd *common.WSCmd) {
	//logger.DebugByName(global.LogFileWSCName, "%v: receive the pong: %v", c.LogPrefix(), cmd.Data)
}
