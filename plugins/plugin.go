package plugins

import (
	"fmt"
	"github.com/vincent78/butil/model"
	"plugin"
)

var PluginMap map[string]PluginModel

type PluginModel struct {
	Path   string
	Plugin *plugin.Plugin
}

func init() {
	PluginMap = make(map[string]PluginModel)
}

func LoadPlugin(name, filePath string) {
	if p, err := plugin.Open(filePath); err != nil {
		panic(err)
	} else {
		PluginMap[name] = PluginModel{
			Path:   filePath,
			Plugin: p,
		}
	}
}

func RunFuncWithArg0(name, funcName string) model.RespModel {
	if p, exist := PluginMap[name]; !exist {
		return model.FailureRespWithStr(500, "not found plugin: %v", name)
	} else {
		if f, err := p.Plugin.Lookup(funcName); err != nil {
			return model.FailureRespWithStr(500, "not found func[%v] in plugin[%v]", funcName, name)
		} else {
			f.(func())()
			return model.SuccessResp("")
		}
	}
}

// 单个参数，单个返回
func RunFuncWithArg[I any, O any](name, funcName string, i I) model.ResultModel[O] {
	if p, exist := PluginMap[name]; !exist {
		return model.FailureResultWithStr[O](500, "not found plugin: %v", name)
	} else {
		if f, err := p.Plugin.Lookup(funcName); err != nil {
			return model.FailureResultWithStr[O](500, "not found func[%v] in plugin[%v]", funcName, name)
		} else {
			o := f.(func(I) O)(i)
			return model.SuccessResult(o)
		}
	}
}

func GetObjFromPlugin(name, objName string) model.RespModel {
	if p, exist := PluginMap[name]; !exist {
		return model.FailureRespWithStr(500, "not found plugin: %v", name)
	} else {
		if f, err := p.Plugin.Lookup(objName); err != nil {
			return model.FailureRespWithError(500, fmt.Errorf("found the obj[%v] in plugin[%v] error: %v", objName, name, err))
		} else {
			return model.SuccessResp(f)
		}
	}
}
