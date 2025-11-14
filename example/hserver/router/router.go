package router

import (
	"github.com/tonly18/xws/example/hserver/controller"
	"net/http"
)

func InitRouter() *http.ServeMux {
	handler := http.NewServeMux()

	//public
	handler.HandleFunc("/v1/public", controller.WrapHandle(controller.PublicController))

	//monitor
	//handler.HandleFunc("/v1/m/goroutine", controller.MonitorGoroutineController)
	//handler.HandleFunc("/v1/m/memory", controller.MonitorMemoryController)
	//test
	//handler.HandleFunc("/v1/test", controller.WrapHandle(controller.TestController))

	//return
	return handler
}
