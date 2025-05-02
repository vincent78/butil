package server

type WSServerConfig struct {
	PingInterval int    `json:"ping,omitempty"` // 心跳的间隔时间
	Version      string `json:"version"`        // 当前版本号
	InfoPath     string `json:"infoPath"`       // the http path of ws infos
}

func NewWSServerConfig() *WSServerConfig {
	return &WSServerConfig{
		PingInterval: 0,
		Version:      "0",
		InfoPath:     "info",
	}
}
