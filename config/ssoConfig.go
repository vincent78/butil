package config

type SsoConfig struct {
	EndPoint   string `json:"endPoint"`   // SSO访问域名
	BucketName string `json:"bucketName"` // SSO中bucketName
	AccessKey  string `json:"accessKey"`  // SSO的用户密码
	AccessId   string `json:"accessId"`   // SSO的用户名
	UseSSL     bool   `json:"useSSL"`     // SSO是否使用SSL
}

func NewSsoConfig() SsoConfig {
	return SsoConfig{}
}
