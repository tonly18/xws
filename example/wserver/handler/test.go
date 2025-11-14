package handler

import (
	"fmt"
	"github.com/tonly18/xws/core/ziface"
	"github.com/tonly18/xws/core/znet"
)

// TestRouter Struct
type TestRouter struct {
	znet.BaseRouter
}

func (h *TestRouter) Handle(request ziface.IRequest) error {
	fmt.Println("handle test...")

	//message
	fmt.Println("data+++++++msg.id: ", request.GetMsgID())
	//fmt.Println("data+++++++msg.data: ", string(request.GetData()))

	//request.GetConnection().SendMsg(201, []byte(`this is a test message from the server!!!`))
	//request.GetConnection().SendByteMsg(201, []byte(`this is a test message from the server!!!`))

	//tcpServer := request.GetConnection().GetServer()
	//fmt.Println("server-hc::::::", tcpServer.GetHeartBeat())
	//
	////panic("test panic")
	//
	//conn := request.GetConnection()
	//fmt.Println("conn-hc::::::", conn.GetHeartBeat())

	fmt.Println("request.GetMsgID()::::::", request.GetMsgID())
	fmt.Println("request.GetData():::::::", string(request.GetData()))

	//return fmt.Errorf("handler is error")
	return PushMessage(request, 101, 0, []byte("a-2023111711111-c"))
}
