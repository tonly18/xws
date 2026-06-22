package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/extra-time-zone/xws/core/zutils"
	httpserver "github.com/extra-time-zone/xws/example/hserver"
	"github.com/extra-time-zone/xws/example/pkg/global"
	"github.com/extra-time-zone/xws/example/service"
	wserver "github.com/extra-time-zone/xws/example/wserver"
)

func main() {
	global.ConfigFile = flag.String("config", "./conf/config_local.toml", "configuration file path")
	flag.Parse()

	//init service
	service.Init()

	//创建监听退出chan,监听指定信号 ctrl+c kill
	sig := zutils.NewSignal()

	/**********************************************************
	 * WS SERVER
	 **********************************************************/
	go func() {
		fmt.Println("[WS SERVER] STARTING UP")
		wserver.StartWsServer(sig)
	}()

	/**********************************************************
	 * HTTP SERVER
	 **********************************************************/
	go func() {
		fmt.Println("[HTTP SERVER] STARTING UP")
		httpserver.StartHttpServer(sig)
	}()

	//Block: wait for signal
	if err := sig.Waiter(); err == nil {
		sig.Cannel(time.Second * 5)
	}

	//Finish
	fmt.Println("[All Services have been Stopped]")
}
