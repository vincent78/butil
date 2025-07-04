package config

import (
	"github.com/vincent78/butil/utils/fileUtil"
	"testing"
)

type CmdConfig struct {
	App    AppConfig    `yaml:"example"`
	Server ServerConfig `yaml:"server"`
	Logger LoggerConfig `yaml:"logger"`
}

func TestParseConfig(t *testing.T) {
	path := fileUtil.GetCurrentProjectPath()
	file := fileUtil.Join(path, "config", "test.yaml")
	conf, err := ParseConfig(file, &CmdConfig{})
	if err != nil {
		t.Errorf("parse the config error: %v", err.Error())
	} else {
		t.Logf("the config : %+v", conf)
	}
}
