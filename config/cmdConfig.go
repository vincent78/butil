package config

type CmdConfig struct {
	App    AppConfig               `yaml:"example"`
	Server ServerConfig            `yaml:"server"`
	Logger map[string]LoggerConfig `yaml:"logger"`
}
