package config

type LoggerConfig struct {
	Home   string `yaml:"home" json:"home"`
	Prefix string `yaml:"prefix" json:"prefix"`
	Level  int    `yaml:"level" json:"level"` // 0 DEBUG,1 INFO,2 WARN,3 ERROR
}

func NewLoggerConfig(prefix string) LoggerConfig {
	return LoggerConfig{
		Home:   "./logs",
		Prefix: prefix,
		Level:  0,
	}
}
