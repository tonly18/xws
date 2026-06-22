package tcpserver

import (
	"fmt"

	"github.com/extra-time-zone/xws/core/logger"
	"github.com/extra-time-zone/xws/core/zconf"
	"github.com/extra-time-zone/xws/core/znet"
	"github.com/extra-time-zone/xws/core/zutils"
	"github.com/extra-time-zone/xws/example/pkg/global"
	"github.com/extra-time-zone/xws/example/sconf"
	"github.com/extra-time-zone/xws/example/wserver/hook"
	"github.com/extra-time-zone/xws/example/wserver/router"
)

func StartWsServer(sig *zutils.Signal) {
	// config
	zconf.Init(&sconf.Config.WsServer)
	logger.Init()

	//0 websocket
	wsServer := global.SetWsServer(znet.NewServer())

	//1 添加router
	router.InitRouter(wsServer)

	//2 Hook
	wsServer.SetOnConnStart(hook.OnConnStartFunc)
	wsServer.SetOnConnStop(hook.OnConnStopFunc)

	//3 启动ws服务
	go func() {
		wsServer.Serve()
	}()

	//4 信号
	select {
	case <-sig.Context().Done():
		fmt.Println("[WS SERVER] SERVER IS STOPPING")
		wsServer.Stop()
	}

	fmt.Println("[WS SERVER] SERVER IS STOPPED")
}
