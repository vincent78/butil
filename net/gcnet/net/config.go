package gcnet

type NetConfig struct {
	Host          string `json:"host" yaml:"host"`
	Port          int    `json:"port" yaml:"port"`
	Protocol      string `json:"protocol" yaml:"protocol"`           // TCP,UDP,UNIX
	Grpc          bool   `json:"grpc" yaml:"grpc"`                   // 是否使用GRPC
	Version       string `json:"version" yaml:"version"`             // 版本
	ReadBuffSize  int    `json:"readBuffSize" yaml:"readBuffSize"`   // 读取的缓冲数组长度
	WriteBuffSize int    `json:"writeBuffSize" yaml:"writeBuffSize"` // 读取的缓冲数组长度
}

func NewNetConfig(host string, port int) *NetConfig {
	return &NetConfig{
		Host:          host,
		Port:          port,
		Protocol:      "tcp",
		Grpc:          false,
		Version:       "0",
		ReadBuffSize:  512,
		WriteBuffSize: 512,
	}
}
