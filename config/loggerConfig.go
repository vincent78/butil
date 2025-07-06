package config

type LoggerConfig struct {
	Format        string        `yaml:"format" json:"format"`
	IsSave        bool          `yaml:"isSave" json:"isSave"`
	Level         string        `yaml:"level" json:"level"`
	DisableCaller bool          `yaml:"disableCaller" json:"disableCaller"`
	LogFileConfig LogFileConfig `yaml:"logFileConfig" json:"logFileConfig"`
}

func NewLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Format:        "console",
		Level:         "debug",
		IsSave:        false,
		DisableCaller: true,
		LogFileConfig: LogFileConfig{
			Filename:      "out.log",
			MaxSize:       20,
			MaxBackups:    50,
			MaxAge:        15,
			IsCompression: true,
		},
	}
}

type Logger struct {
	Format        string        `yaml:"format" json:"format"`
	IsSave        bool          `yaml:"isSave" json:"isSave"`
	Level         string        `yaml:"level" json:"level"`
	LogFileConfig LogFileConfig `yaml:"logFileConfig" json:"logFileConfig"`
}

type LogFileConfig struct {
	Filename      string `yaml:"filename" json:"filename"`
	IsCompression bool   `yaml:"isCompression" json:"isCompression"`
	MaxAge        int    `yaml:"maxAge" json:"maxAge"`
	MaxBackups    int    `yaml:"maxBackups" json:"maxBackups"`
	MaxSize       int    `yaml:"maxSize" json:"maxSize"`
}
