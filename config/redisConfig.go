package config

import "fmt"

type RedisConfig struct {
	Host     string `yaml:"host" json:"host"`         // 服务器地址
	Port     int    `yaml:"port" json:"port"`         // 服务端口号
	Db       int    `yaml:"db" json:"db"`             // 数据库名称
	UserName string `yaml:"userName" json:"userName"` // 用户名
	Passwd   string `yaml:"passwd" json:"passwd"`     //数据库密码
}

func NewRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     "127.0.0.1",
		Port:     6379,
		Db:       0,
		UserName: "",
		Passwd:   "",
	}
}

func (conf *RedisConfig) UrlStr() string {

	/*
		转义：
		+ 转义后 %2B
		空格 转义后 %20
		/ 转义后 %2F
		? 转义后 %3F
		% 转义后 %25
		# 转义后 %23
		& 转义后 %26
		= 转义后 %3D

		redis.address=redis://:redis#123@127.0.0.1:6379/0

		redis.address=redis://:redis%23123@127.0.0.1:6379/0
	*/

	if conf.Port == 0 {
		conf.Port = 6379
	}

	if conf.Host == "" {
		conf.Host = "127.0.0.1"
	}

	return fmt.Sprintf("redis://:%v@%v:%v/%v",
		conf.Passwd,
		conf.Host,
		conf.Port,
		conf.Db,
	)
}
