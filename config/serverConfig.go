package config

type ServerConfig struct {
	Http *ServerDetailConfig `yaml:"http" json:"http,omitempty"`
	Ws   *ServerDetailConfig `yaml:"ws" json:"ws,omitempty"`
	Tcp  *ServerDetailConfig `yaml:"tcp" json:"tcp,omitempty"`
	Udp  *ServerDetailConfig `yaml:"udp" json:"udp,omitempty"`
}

type ServerDetailConfig struct {
	Host string `yaml:"host" json:"host"`
	Port int    `yaml:"port" json:"port"`
	SSL  bool   `yaml:"ssl" json:"ssl,omitempty"`
	Ext  string `yaml:"ext" json:"ext,omitempty"`
}
