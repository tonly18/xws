package znet

import (
	"sync"
)

const maxSize = 4

var packetPool = &sync.Pool{
	New: func() any {
		return make([]byte, maxSize)
	},
}

func packetGetFromPool() []byte {
	return packetPool.Get().([]byte)
}

func packetPutToPool(buffer []byte) {
	if cap(buffer) <= maxSize {
		clear(buffer)
		packetPool.Put(buffer)
	}
}
