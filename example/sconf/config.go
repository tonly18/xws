package sconf

import (
	"fmt"
	"github.com/tonly18/xws/core/zconf"

	"github.com/spf13/viper"
)

type config struct {
	Http     HttpConfig
	WsServer zconf.ZConfig
}

type HttpConfig struct {
	Host string
	Port int
}

var Config *config = &config{}

// Init 初始化配置信息
func Init(file *string) error {
	if err := parseConfigFromToml(file); err != nil {
		panic(fmt.Sprintf("config init from file error:%v", err))
	}

	//fmt.Printf("config:%+v\n", Config)

	return nil
}

// 获取配置文件并解析到指定的变量
func parseConfigFromToml(file *string) error {
	viper.SetConfigFile(*file)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("FileName: %s, Error: %v", file, err)
	}

	//parse
	return viper.Unmarshal(Config)
}
