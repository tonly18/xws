package handler

import (
	"github.com/tonly18/xws/core/ziface"
)

// PushMessage 向客户端推送消息
func PushMessage(req ziface.IRequest, cmd, code uint32, data []byte) error {
	//downMsg := pack.NewMessageDown(cmd, code, data)
	//dp := pack.NewDataPackDown()
	//msgPack, err := dp.Pack(downMsg)
	//if err != nil {
	//	return fmt.Errorf(`base push message unpack error:%v`, err)
	//}
	//if err := req.GetConnection().SendBuffMsg(msgPack); err != nil {
	//	return fmt.Errorf(`base push message conn send error:%v`, err)
	//}

	//return
	return nil
}
