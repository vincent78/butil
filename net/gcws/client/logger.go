package client

import "fmt"

func (c *WSClient) NormalLogger(msg string, args ...any) string {
	str := fmt.Sprintf(msg, args...)
	return fmt.Sprintf("%v: %v ", c.LogPrefix(), str)
}

func (c *WSClient) LogPrefix() string {
	if len(c.ClientAddr) > 0 {
		return fmt.Sprintf("==== [%v]", c.ClientAddr)
	} else {
		return fmt.Sprintf("==== [%v]", ":")
	}
}

func (c *WSClient) LogPrefixIn() string {
	return ">>>>"
}

func (c *WSClient) LogPrefixOut() string {
	return "<<<<"
}
