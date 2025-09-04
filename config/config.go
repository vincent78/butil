package config

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
	"github.com/vincent78/butil/utils/fileUtil"
	"github.com/vincent78/butil/utils/strUtil"
)

/*

# 配置文件加载的优先级
1、命令行参数
2、配置文件（手动指定）
3、环境变量
4、默认值（默认的配置文件或原始值)
*/

type ChangedHandler func(key string, value interface{})

var (
	globalConfig interface{}
)

// ParseConfig conf 必须传入指针类型
func ParseConfig[T any](filepath string, conf T) (T, error) {
	//configMap[filepath] = val
	fmt.Println("configPath:", filepath)

	if filepath == "" {
		return conf, nil
	}

	//导入配置文件
	viper.SetConfigType("yaml")
	viper.SetConfigFile(fileUtil.GetAbs(filepath))
	//读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
		return conf, err
	}
	//将配置文件读到结构体中
	err = viper.Unmarshal(&conf)
	if err != nil {
		fmt.Println(err.Error())
		return conf, err
	}

	confBytes, err := json.Marshal(conf)
	fmt.Print(strUtil.Bytes2String(confBytes))
	fmt.Print("\n\n\n")
	globalConfig = conf
	return conf, nil
}

func SubscribeConfigChanged(key string, handler ChangedHandler) {

}
