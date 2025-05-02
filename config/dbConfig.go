package config

type DBInterface interface {
	// UrlStr 获取数据库连接字符串
	UrlStr() string
}

type DBConfig struct {
	Type   string `yaml:"type" json:"type"`     // 数据库类型
	Host   string `yaml:"host" json:"host"`     // 服务器地址
	Port   int    `yaml:"port" json:"port"`     // 服务端口号
	Db     string `yaml:"db" json:"db"`         // 数据库名称
	User   string `yaml:"user" json:"user"`     // 数据库用户名
	Passwd string `yaml:"passwd" json:"passwd"` //数据库密码
}
