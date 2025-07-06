package config

type CmdConfig struct {
	App      AppConfig               `yaml:"app" json:"app"`
	Server   ServerConfig            `yaml:"server" json:"server"`
	Logger   map[string]LoggerConfig `yaml:"logger" json:"logger"`
	Database Database                `yaml:"database" json:"database"`
	Redis    Redis                   `yaml:"redis" json:"redis"`
	Chain    Chain                   `yaml:"chain" json:"chain"`
}

type Database struct {
	Type string `json:"type" yaml:"type"`
	Dev  Mysql  `json:"dev" yaml:"dev"`
}

type Chain struct {
	Tron Tron `json:"tron" yaml:"tron"`
}

type Tron struct {
	RunType   string     `json:"runType" yaml:"runType"`
	Callbacks []Callback `json:"callback" yaml:"callback"`
	Apis      Apis       `json:"apis" yaml:"apis"`
	Accounts  []Account  `json:"accounts" yaml:"accounts"`
	Url       string     `json:"url" yaml:"url"`
	Contracts []Contract `json:"contracts" yaml:"contracts"`
	Ext       string     `json:"extra" yaml:"extra"`
}

type Callback struct {
	Name string `yaml:"name" json:"name"`
	Url  string `yaml:"url" json:"url"`
	Ext  string `yaml:"ext" json:"ext"`
}
type Apis struct {
	TronGrip []TronGripApi `json:"tron" yaml:"tron"`
}

type Account struct {
	Name    string `json:"name" yaml:"name"`
	Address string `json:"address" yaml:"address"`
	PriKey  string `json:"priKey" yaml:"priKey"`
}

type TronGripApi struct {
	Name   string `yaml:"name" json:"name"`
	Key    string `yaml:"key" json:"key"`
	Token  string `yaml:"token" json:"token"`
	Weight int    `yaml:"weight" json:"weight"`
}

type Contract struct {
	Name    string `yaml:"name" json:"name"`
	Address string `yaml:"address" json:"address"`
	Decimal int    `yaml:"decimal" json:"decimal"`
}
