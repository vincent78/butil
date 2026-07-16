package logger4

// LoggerConfig is the lightweight configuration owned by logger4.
//
// Keep logger4 independent from the global config package: callers that read
// config files can map their own config struct to this one at the application
// boundary.
type LoggerConfig struct {
	Format          string        `yaml:"format" json:"format"`
	IsSave          bool          `yaml:"isSave" json:"isSave"`
	Level           string        `yaml:"level" json:"level"`
	DisableCaller   bool          `yaml:"disableCaller" json:"disableCaller"`
	IsLocalTime     bool          `yaml:"isLocalTime" json:"isLocalTime"`
	CallerSkip      int           `yaml:"callerSkip" json:"callerSkip"`
	StacktraceLevel string        `yaml:"stacktraceLevel" json:"stacktraceLevel"`
	LogFileConfig   LogFileConfig `yaml:"logFileConfig" json:"logFileConfig"`
}

type LogFileConfig struct {
	Filename      string `yaml:"filename" json:"filename"`
	IsCompression bool   `yaml:"isCompression" json:"isCompression"`
	MaxAge        int    `yaml:"maxAge" json:"maxAge"`
	MaxBackups    int    `yaml:"maxBackups" json:"maxBackups"`
	MaxSize       int    `yaml:"maxSize" json:"maxSize"`
	IsLocalTime   bool   `yaml:"isLocalTime" json:"isLocalTime"`
}

func NewLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Format:          "console",
		Level:           "debug",
		IsSave:          false,
		DisableCaller:   true,
		CallerSkip:      1,
		StacktraceLevel: "error",
		LogFileConfig: LogFileConfig{
			Filename:      "out.log",
			MaxSize:       256,
			MaxBackups:    15,
			MaxAge:        20,
			IsCompression: true,
			IsLocalTime:   false,
		},
	}
}
