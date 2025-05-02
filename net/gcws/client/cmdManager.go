package client

import (
	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/logger"
	"github.com/vincent78/butil/net/gcws/common"
)

func DoCmdAction(c *WSClient, cmd *common.WSCmd) {
	if cmd.Cmd == "" && cmd.PID != "" {
		logger.DebugByName(global.LogFileWSSName, c.NormalLogger("receive the resp without handler"))
	} else {
		cmdKey := cmd.HandlerKey()
		if handler, ok := c.cmdHandlers[cmdKey]; ok && handler != nil {
			handler.doAction(c, cmd)
		} else {
			errModel := common.ErrorWSSNoCmdName(cmdKey)
			logger.ErrorByName(global.LogFileWSCName, c.NormalLogger(errModel.Message))
			c.Resp(cmd.FailureByErrModel(errModel))
		}
	}
}
