package config

type AppConfig struct {
	Name  string `yaml:"name" json:"name"`
	Debug bool   `yaml:"debug" json:"debug"`
	Proxy string `yaml:"proxy" json:"proxy"`
}
