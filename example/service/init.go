package service

import (
	"github.com/tonly18/xws/example/pkg/global"
	"github.com/tonly18/xws/example/sconf"
)

func Init() {
	//init config
	sconf.Init(global.ConfigFile)
}
