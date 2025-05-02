package server

import (
	"gitee.com/vincent78/gcutil/net/gcws/common"
	"gitee.com/vincent78/gcutil/utils/timeUtil"
)

type PingCmdHandler struct {
}

func (h *PingCmdHandler) DoAction(c *WSServer, cmd *common.WSCmd) {
	//logger.DebugByName(global.LogFileWSSName, c.NormalLogger("%v", cmd.Data))
	resp := common.NewWSCmdByCategory(common.CmdCategorySys, common.CmdNamePong)
	resp.Data = timeUtil.NowStr()
	resp.PID = cmd.ID
	c.Send(resp)
}
