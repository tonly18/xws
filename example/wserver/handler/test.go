package handler

import (
	"fmt"

	"github.com/extra-time-zone/xws/core/ziface"
	"github.com/extra-time-zone/xws/core/znet"
	"github.com/extra-time-zone/xws/example/pkg/global"
)

// TestRouter Struct
type TestRouter struct {
	znet.BaseRouter
}

func (h *TestRouter) Handle(req ziface.IRequest) error {
	//fmt.Println("handle test...")

	//message
	//fmt.Println("data+++++++msg.id: ", req.GetMsgID())
	//fmt.Println("data+++++++msg.data: ", string(req.GetData()))

	//req.GetConnection().SendMsg(201, []byte(`this is a test message from the server!!!`))
	//req.GetConnection().SendByteMsg(201, []byte(`this is a test message from the server!!!`))

	//tcpServer := req.GetConnection().GetServer()
	//fmt.Println("server-hc::::::", tcpServer.GetHeartBeat())

	//conn := req.GetConnection()
	//fmt.Println("conn-hc::::::", conn.GetHeartBeat())

	fmt.Printf("---TestRouter MsgID:%d Data:%s\n", req.GetMsgID(), string(req.GetData()))

	//return fmt.Errorf("handler is error")
	return PushMessage(req, global.DownCmdTest, []byte("test-message-20260302"))
}
