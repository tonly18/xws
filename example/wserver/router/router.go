package router

import (
	"github.com/extra-time-zone/xws/core/ziface"
	"github.com/extra-time-zone/xws/example/pkg/global"
	"github.com/extra-time-zone/xws/example/wserver/handler"
)

func InitRouter(s ziface.IServer) {
	//test handler
	s.AddRouter(global.UpCmdTest, &handler.TestRouter{})
}
