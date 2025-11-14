package zconf

import "time"

type ZConfig struct {
	ENV      string
	Host     string //当前服务器主机IP
	Port     int    //当前服务器主机监听端口号
	Path     string
	Name     string //当前服务器名称
	ServerID uint32 //服务器ID

	MaxConn          int    //当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32 //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint32 //SendBuffMsg发送消息的缓冲最大长度
	MaxPacketSize    uint32 //数据包的最大值

	//conn 读写时间
	MaxConnReadTime  int //conn 读时间：单位毫秒
	MaxConnWriteTime int //conn 写时间：单位毫秒

	//https 证书
	CertFile       string
	PrivateKeyFile string

	//最长心跳检测间隔时长,超过改时间间隔,则认为超时
	HeartbeatMax time.Duration
}

var Config *ZConfig

func Init(c *ZConfig) {
	Config = &ZConfig{
		ENV:      LocalEnv,
		Host:     "0.0.0.0",
		Port:     7000,
		Path:     "/ws",
		Name:     "xp-websocket",
		ServerID: 1,

		MaxConn:          10000,
		WorkerPoolSize:   16,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
		MaxPacketSize:    1024,

		MaxConnReadTime:  1800,
		MaxConnWriteTime: 1800,

		HeartbeatMax: time.Second * 30,
	}

	if c.ENV != "" {
		Config.ENV = c.ENV
	}
	if c.Host != "" {
		Config.Host = c.Host
	}
	if c.Port != 0 {
		Config.Port = c.Port
	}
	if c.Path != "" {
		Config.Path = c.Path
	}
	if c.Name != "" {
		Config.Name = c.Name
	}
	if c.ServerID != 0 {
		Config.ServerID = c.ServerID
	}
	if c.MaxConn != 0 {
		Config.MaxConn = c.MaxConn
	}
	if c.WorkerPoolSize != 0 {
		Config.WorkerPoolSize = c.WorkerPoolSize
	}
	if c.MaxWorkerTaskLen != 0 {
		Config.MaxWorkerTaskLen = c.MaxWorkerTaskLen
	}
	if c.MaxMsgChanLen != 0 {
		Config.MaxMsgChanLen = c.MaxMsgChanLen
	}
	if c.MaxPacketSize != 0 {
		Config.MaxPacketSize = c.MaxPacketSize
	}
	if c.MaxConnReadTime != 0 {
		Config.MaxConnReadTime = c.MaxConnReadTime
	}
	if c.MaxConnWriteTime != 0 {
		Config.MaxConnWriteTime = c.MaxConnWriteTime
	}
	if c.CertFile != "" {
		Config.CertFile = c.CertFile
	}
	if c.PrivateKeyFile != "" {
		Config.PrivateKeyFile = c.PrivateKeyFile
	}
	if c.HeartbeatMax != 0 {
		Config.HeartbeatMax = c.HeartbeatMax
	}
}
