package service

import (
	"github.com/extra-time-zone/xws/example/pkg/global"
	"github.com/extra-time-zone/xws/example/sconf"
)

func Init() {
	//init config
	sconf.Init(global.ConfigFile)
}
