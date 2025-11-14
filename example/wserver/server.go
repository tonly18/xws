package tcpserver

import (
	"fmt"
	"github.com/tonly18/xws/core/logger"
	"github.com/tonly18/xws/core/zconf"
	"github.com/tonly18/xws/core/znet"
	"github.com/tonly18/xws/core/zutils"
	"github.com/tonly18/xws/example/global"
	"github.com/tonly18/xws/example/sconf"
	"github.com/tonly18/xws/example/wserver/hook"
	"github.com/tonly18/xws/example/wserver/router"
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
	case <-sig.GetCtx().Done():
		fmt.Println("[WS SERVER] SERVER IS STOPPING")
		wsServer.Stop()
	}

	fmt.Println("[WS SERVER] SERVER IS STOPPED")
}
