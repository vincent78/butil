package config

import (
	"github.com/vincent78/butil/utils/fileUtil"
	"testing"
)

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

func TestParseConfig2(t *testing.T) {
	path := fileUtil.GetCurrentProjectPath()
	file := fileUtil.Join(path, "config", "test.yaml")
	conf := &CmdConfig{}
	err := Parse(file, conf)
	if err != nil {
		t.Errorf("parse the config error: %v", err.Error())
	}
	t.Logf("the config : %+v", conf)
}
