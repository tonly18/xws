package znet

import (
	"context"
	"sync"

	"github.com/extra-time-zone/xws/core/ziface"
)

var requestPool = &sync.Pool{
	New: func() any {
		return requestAllocate()
	},
}

func requestAllocate() ziface.IRequest {
	return NewRequest(nil, nil)
}

func requestGetFromPool(conn ziface.IConnection, msg ziface.IMessage) ziface.IRequest {
	req := requestPool.Get().(*Request)
	req.conn = conn
	req.msg = msg
	if req.ctx == nil {
		req.ctx = context.Background()
	}

	return req
}

func requestPutToPool(req ziface.IRequest) {
	req.Reset()
	requestPool.Put(req)
}
