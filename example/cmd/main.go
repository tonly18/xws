package main

import (
	"flag"
	"fmt"
	"github.com/tonly18/xws/core/zutils"
	"github.com/tonly18/xws/example/global"
	httpserver "github.com/tonly18/xws/example/hserver"
	"github.com/tonly18/xws/example/service"
	wserver "github.com/tonly18/xws/example/wserver"
)

func main() {
	global.ConfigFile = flag.String("config", "../conf/config_local.toml", "configuration file path")
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
		sig.Cannel()
	}

	//Finish
	fmt.Println("[All Services have been Stopped]")
}
