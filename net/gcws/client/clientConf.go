package client

type WSClientConfig struct {
	PingInterval      int64 // 间隔 （秒）
	ConnRetryInterval int64 // 重次的间隙
	ConnMaxNum        int   // 重试的次数
}

func NewWSClientConfig() *WSClientConfig {
	return &WSClientConfig{
		PingInterval:      90,
		ConnRetryInterval: 5,
		ConnMaxNum:        30,
	}
}
