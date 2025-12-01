package config

type AppConfig struct {
	Name  string `yaml:"name" json:"name"`
	Node  string `yaml:"node" json:"node"`
	Debug bool   `yaml:"debug" json:"debug"`
	Proxy string `yaml:"proxy" json:"proxy"`
}
