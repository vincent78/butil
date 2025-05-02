package server

import (
	"github.com/vincent78/butil/net/gcws/common"
)

//func ParseCmdInput[T any](source interface{}, obj T) *model.ErrorModel {
//	if source == nil {
//		return nil
//	}
//	bytes, err := json.Marshal(source)
//	if err != nil {
//		return model.ErrorBaseNormal(err.Error())
//	}
//	err = json.Unmarshal(bytes, obj)
//	if err != nil {
//		return model.ErrorBaseNormal(err.Error())
//	}
//	return nil
//}

type WSServerCmdHandlerInterface interface {
	DoAction(c *WSServer, cmd *common.WSCmd)
}

func (c *WSServer) initCmdHandler() {
	c.RegisterCmdHandler(common.CmdCategorySys, common.CmdNameLogin, &LoginCmdHandler{})
	c.RegisterCmdHandler(common.CmdCategorySys, common.CmdNameLogout, &LogoutCmdHandler{})
	c.RegisterCmdHandler(common.CmdCategorySys, common.CmdNameClose, &CloseCmdHandler{})
	c.RegisterCmdHandler(common.CmdCategorySys, common.CmdNamePing, &PingCmdHandler{})
}
