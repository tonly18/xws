package hook

import (
	"fmt"

	"github.com/extra-time-zone/xws/core/ziface"
)

func OnConnStartFunc(conn ziface.IConnection) {
	fmt.Printf("[WebSocket Server] r.Header.Accept-Language:%+v\n", conn.GetRequest().Header.Get("Accept-Language"))

	conn.SetProperty("Accept-Language", conn.GetRequest().Header.Get("Accept-Language"))

	fmt.Println("[WebSocket Server] OnConnStartFunc, Accept-Language:", conn.GetProperty("Accept-Language"))
}

func OnConnStopFunc(conn ziface.IConnection) {
	fmt.Println("[WebSocket Server] OnConnStopFunc connId:", conn.GetConnID())
}
