package server

import (
	"fmt"

	"github.com/vincent78/butil/global"
	"github.com/vincent78/butil/logger/logger1"
)

func (c *WSServer) NormalLogger(msg string, args ...any) string {
	str := fmt.Sprintf(msg, args...)
	return c.loggerMsg("====", str)
}

func (c *WSServer) ILog(msg string, args ...any) {
	logger1.InfoByName(global.LogFileWSSName, c.InLogger(msg, args...))
}

func (c *WSServer) InLogger(msg string, args ...any) string {
	str := fmt.Sprintf(msg, args...)
	return c.loggerMsg(">>>>", str)
}

func (c *WSServer) OLog(msg string, args ...any) {
	logger1.InfoByName(global.LogFileWSSName, c.OutLogger(msg, args...))
}

func (c *WSServer) OutLogger(msg string, args ...any) string {
	str := fmt.Sprintf(msg, args...)
	return c.loggerMsg("<<<<", str)
}

func (c *WSServer) loggerMsg(prefix, str string) string {
	if c.ID != "" {
		return fmt.Sprintf("%v [%v][%v]: %v", prefix, c.RemoteIP, c.ID, str)
	} else {
		return fmt.Sprintf("%v [%v]: %v", prefix, c.RemoteIP, str)
	}
}
