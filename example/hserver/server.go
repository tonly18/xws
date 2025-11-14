package httpserver

import (
	"errors"
	"fmt"
	"github.com/tonly18/xws/core/zserver"
	"github.com/tonly18/xws/core/zutils"
	"github.com/tonly18/xws/example/hserver/router"
	"github.com/tonly18/xws/example/sconf"
	"net/http"
)

func StartHttpServer(sig *zutils.Signal) {
	//0 create http server
	config := zserver.HttpServerConfig{
		IP:      sconf.Config.Http.Host,
		Port:    sconf.Config.Http.Port,
		Handler: router.InitRouter(),
	}
	httpserver := zserver.NewHttpServer(&config)

	//1 listen
	go func() {
		if err := httpserver.Start(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(fmt.Sprintf("[HTTP SERVER] LISTEN ERROR: %v", err))
			}
		}
	}()

	//2 signal
	select {
	case <-sig.GetCtx().Done():
		fmt.Println("[HTTP SERVER] SERVER IS STOPPING")
		httpserver.Stop()
	}

	fmt.Println("[HTTP SERVER] SERVER IS STOPPED")
}
