package sconf

import (
	"fmt"

	"github.com/tonly18/xws/core/zconf"

	"github.com/spf13/viper"
)

type config struct {
	Http     HttpConfig
	WsServer zconf.ZConfig
	MySql    []*MySqlConfig
	Redis    *RedisConfig
}

type HttpConfig struct {
	Host string
	Port int
}

type MySqlConfig struct {
	Role         string         `json:"role"`
	Host         string         `json:"host"`
	Port         int            `json:"port"`
	Dbname       string         `json:"dbname"`
	Username     string         `json:"username"`
	Password     string         `json:"password"`
	Charset      string         `json:"charset"`
	Collation    string         `json:"collation"`
	MaxIdleConns int            `json:"max_idle_conns"`
	MaxOpenConns int            `json:"max_open_conns"`
	MaxLifetime  int            `json:"max_lifetime"`
	MaxIdleTime  int            `json:"max_idle_time"`
	Slave        []*MySqlConfig `json:"slave"`
}

type RedisConfig struct {
	Host         []string `json:"host"`
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	ServerName   string   `json:"server_name"`
	MaxIdleConns int      `json:"max_idle_conns"` //最大空闲链接数
	MinIdleConns int      `json:"min_idle_conns"` //最小空闲链接数
	PoolSize     int      `json:"pool_size"`      //链接池最大链接数
	Prefix       string   `json:"prefix"`         //key前缀
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
