package ziface

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
)

type IConnection interface {
	Start()                   //启动连接，让当前连接开始工作
	Stop()                    //停止连接，结束当前连接状态M
	Context() context.Context //返回ctx，用于用户自定义的go程获取连接退出状态

	GetConnection() *websocket.Conn //从当前连接获取原始的socket Conn
	SetConnID(uint64)               //设置当前连接ID
	GetConnID() uint64              //获取当前连接ID
	GetConnMgr() IConnManager       //获取connection管理器
	GetMsgHandler() IMsgHandle      //获取消息处理器
	GetRemoteAddr() string          //获取远程客户端地址信息
	GetLocalAddr() string           //获取服务端地址信息
	GetClientIP() string            //获取客户端IP
	GetName() string

	Send([]byte) error
	SendBuffMsg(uint32, []byte) error //直接将二进制流发送给远程的客户端(有缓冲)

	SetProperty(any, any) //设置链接属性
	GetProperty(any) any  //获取链接属性
	RemoveProperty(any)   //移除链接属性

	IsAlive() bool             //判断当前连接是否存活
	GetRequest() *http.Request //获取请求对象
}
