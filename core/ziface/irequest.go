package ziface

import (
	"context"
	"time"
)

type HandleStep int8

type IRequest interface {
	GetConnection() IConnection //获取请求连接信息
	GetData() []byte            //获取请求消息的数据
	GetMsgID() uint32           //获取请求的消息ID
	BindRouter(IRouter)         //绑定这次请求由哪个路由处理
	Call() error                //转进到下一个处理器开始执行,但是调用此方法的函数会根据先后顺序逆序执行
	//Abort()                   //终止处理函数的运行,但调用此方法的函数会执行完毕

	SetAargs(string, any)
	GetAargs(string) any
	GetTraceId() string
	Reset()

	GetCtx() context.Context
	Deadline() (time.Time, bool)
	Done() <-chan struct{}
	Err() error
	Value(any) any
}
