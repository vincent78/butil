package client

import (
	"gitee.com/vincent78/gcutil/global"
	"gitee.com/vincent78/gcutil/logger"
	"gitee.com/vincent78/gcutil/net/gcws/common"
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
