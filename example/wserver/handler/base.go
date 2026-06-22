package handler

import (
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/extra-time-zone/xws/core/ziface"
)

// PushMessage 向客户端推送消息
func PushMessage(req ziface.IRequest, cmd uint32, data any) error {
	var dataBytes []byte

	switch v := data.(type) {
	case []byte:
		dataBytes = v
	default:
		dataByte, err := sonic.Marshal(data)
		if err != nil {
			return fmt.Errorf(`base push message json.Marshal error:%v`, err)
		}
		dataBytes = dataByte
	}

	if err := req.GetConnection().SendBuffMsg(cmd, dataBytes); err != nil {
		return fmt.Errorf(`base push message error:%v`, err)
	}

	return nil
}
