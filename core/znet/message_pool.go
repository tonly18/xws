package znet

import (
	"github.com/tonly18/xws/core/ziface"
	"sync"
)

var messagePool = &sync.Pool{
	New: func() any {
		return allocateMessage()
	},
}

func allocateMessage() ziface.IMessage {
	return &Message{
		Cmd:  0,
		Data: make([]byte, 0, 256),
	}
}

func messageGetFromPool() ziface.IMessage {
	return messagePool.Get().(*Message)
}

func messagePutToPool(msg ziface.IMessage) {
	msg.SetCmd(0)
	msg.SetData(nil)
	messagePool.Put(msg)
}
