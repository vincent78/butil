package client

import (
	"fmt"
	"gitee.com/vincent78/gcutil/net/gcws/common"
)

type WSClientCmdHandlerInterface interface {
	doAction(c *WSClient, cmd *common.WSCmd)
}

func (c *WSClient) initSysCmdHander() {
	c.RegisterCmdHandler(common.CmdCategorySys, common.CmdNameClose, &CloseCMDHandler{})
	c.RegisterCmdHandler(common.CmdCategorySys, common.CmdNamePong, &PongCMDHandler{})
}

func (c *WSClient) RegisterCmdHandler(category, cmd string, handler WSClientCmdHandlerInterface) {
	c.cmdHandlers[fmt.Sprintf("%v:%v", category, cmd)] = handler
}
