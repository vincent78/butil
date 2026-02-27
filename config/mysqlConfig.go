package config

import "fmt"

type MySqlConfig struct {
	DBConfig `json:",inline" yaml:",inline" mapstructure:",squash"`
	Charset  string `json:"charset"`
	TimeZone string `json:"timeZone"`
	Version  string `json:"version"`
}

func NewMySqlConfig() MySqlConfig {
	return MySqlConfig{
		DBConfig: DBConfig{
			Type: "mysql",
			Host: "127.0.0.1",
			Port: 3306,
		},
		Charset:  "utf8mb4",
		TimeZone: "Asia/Shanghai",
		Version:  "1",
	}
}

func (conf *MySqlConfig) UrlStr() string {
	if conf.Charset == "" {
		conf.Charset = "utf8"
	}

	if conf.Port == 0 {
		conf.Port = 3306
	}

	if conf.TimeZone == "" {
		conf.TimeZone = "Asia/Shanghai"
	}

	if conf.Host == "" {
		conf.Host = "127.0.0.1"
	}
	//%s:%s@tcp(%v:%v)/%v&charset=%s&autoReconnect=true&parseTime=True&loc=Local&authSource=%v&authMechanism=SCRAM-SHA-256
	// "%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local&timeout=1000ms"
	if conf.Version == "" || conf.Version == "1" {
		return fmt.Sprintf("%s:%s@tcp(%v:%v)/%v?charset=%s&parseTime=True&loc=Local",
			conf.User,
			conf.Passwd,
			conf.Host,
			conf.Port,
			conf.Db,
			conf.Charset,
		)
	} else {
		return fmt.Sprintf("%s:%s@tcp(%v:%v)/%v&charset=%s&autoReconnect=true&parseTime=True&loc=Local&authSource=%v&authMechanism=SCRAM-SHA-256",
			conf.User,
			conf.Passwd,
			conf.Host,
			conf.Port,
			conf.Db,
			conf.Charset,
			conf.Db,
		)
	}
}

type Mysql struct {
	Dsn             string `yaml:"dsn" json:"dsn"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime" json:"connMaxLifetime"`
	EnableLog       bool   `yaml:"enableLog" json:"enableLog"`
	Logger          string `yaml:"logger" json:"logger"`
	LogLevel        string `yaml:"logLevel" json:"logLevel"`
	TraceIdKey      string `yaml:"traceIdKey" json:"traceIdKey"`
	CacheType       string `yaml:"cacheType" json:"cacheType"`
	MaxIdleConns    int    `yaml:"maxIdleConns" json:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns" json:"maxOpenConns"`
}
