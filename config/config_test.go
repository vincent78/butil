package config

import "testing"

const configFile = "/tmp/test.yaml"

type CmdConfig struct {
	App    AppConfig    `yaml:"example"`
	Server ServerConfig `yaml:"server"`
}

func TestParseConfig(t *testing.T) {
	conf, err := ParseCofnig(configFile, &CmdConfig{})
	if err != nil {
		t.Errorf("parse the config error: %v", err.Error())
	} else {
		t.Logf("the config : %+v", conf)
	}
}
