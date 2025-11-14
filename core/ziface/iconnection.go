package ziface

import (
	"context"

	"github.com/gorilla/websocket"
)

type IConnection interface {
	Start()                   //启动连接，让当前连接开始工作
	Stop()                    //停止连接，结束当前连接状态M
	Context() context.Context //返回ctx，用于用户自定义的go程获取连接退出状态

	GetConnection() *websocket.Conn //从当前连接获取原始的socket Conn
	GetConnID() uint64              //获取当前连接ID
	GetConnMgr() IConnManager       //获取connection管理器
	GetMsgHandler() IMsgHandle      //获取消息处理器
	GetRemoteAddr() string          //获取远程客户端地址信息
	GetLocalAddr() string           //获取服务端地址信息
	GetName() string

	Send([]byte) error
	SendBuffMsg(uint32, []byte) error //直接将二进制流发送给远程的客户端(有缓冲)

	SetProperty(string, any) //设置链接属性
	GetProperty(string) any  //获取链接属性
	RemoveProperty(string)   //移除链接属性

	//GetCreateTime() time.Time //获取conn创建时间
	//GetActivity() time.Time   //获取conn最后活跃时间
	//UpdateActivity()          //更新conn最后活跃时间

	IsAlive() bool //判断当前连接是否存活
}
