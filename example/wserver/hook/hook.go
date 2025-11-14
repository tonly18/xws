package hook

import (
	"fmt"
	"github.com/tonly18/xws/core/ziface"
)

func OnConnStartFunc(conn ziface.IConnection) {
	fmt.Println("[WebSocket Server] OnConnStartFunc")
}

func OnConnStopFunc(conn ziface.IConnection) {
	fmt.Println("[WebSocket Server] OnConnStopFunc")
}
