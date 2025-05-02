package server

import (
	"github.com/vincent78/butil/net/gcws/common"
	"github.com/vincent78/butil/utils/timeUtil"
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
