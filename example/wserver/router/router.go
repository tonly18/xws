package router

import (
	"github.com/tonly18/xws/core/ziface"
	"github.com/tonly18/xws/example/wserver/handler"
)

func InitRouter(s ziface.IServer) {
	//test handler
	s.AddRouter(0, &handler.TestRouter{})

	//proto handler
	//s.AddRouter(global.CMD_UP_PROTO, &handler.PublicRouter{})
	//
	////login handler
	//s.AddRouter(global.CMD_UP_LOGIN, &handler.LoginRouter{})
	//
	////ping handler
	//s.AddRouter(global.CMD_UP_PING, &handler.PingRouter{})
}
