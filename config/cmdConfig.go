package config

type CmdConfig struct {
	App      AppConfig               `yaml:"app" json:"app"`
	Server   ServerConfig            `yaml:"server" json:"server"`
	Logger   map[string]LoggerConfig `yaml:"logger" json:"logger"`
	Database Database                `yaml:"database" json:"database"`
	Redis    Redis                   `yaml:"redis" json:"redis"`
}

type Database struct {
	Type string `json:"type" yaml:"type"`
	Dev  Mysql  `json:"dev" yaml:"dev"`
}
