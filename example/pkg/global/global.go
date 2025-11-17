package global

import (
	"github.com/tonly18/xws/core/ziface"
)

var (
	// ServerEnv 运行环境: local、dev、test、prod
	ServerEnv string

	// ConfigFile 配置文件
	ConfigFile *string
)

// 全局 ws Server
var globalWsServer ziface.IServer

// SetWsServer 获取wsServer
func SetWsServer(wsServer ziface.IServer) ziface.IServer {
	globalWsServer = wsServer

	return globalWsServer
}

// GetWsServer 获取tcpServer
func GetWsServer() ziface.IServer {
	return globalWsServer
}
