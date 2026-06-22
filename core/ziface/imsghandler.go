package ziface

type IMsgHandle interface {
	SendMsgToTaskQueue(IRequest) //将消息交给TaskQueue,由worker进行处理
	DoMsgHandler(IRequest)       //马上以非阻塞方式处理消息
	AddRouter(uint32, IRouter)   //为消息添加具体的处理逻辑
	StartWorkerPool()            //启动worker工作池
}
