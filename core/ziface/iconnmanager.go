package ziface

type IConnManager interface {
	Add(IConnection)                 //添加链接
	Get(uint64) (IConnection, error) //利用ConnID获取链接
	GetConns([]uint64) []IConnection //利用ConnID获取多个链接
	GetAllConns() []IConnection      //获取所有Conn链接
	RangeConn(fn func(IConnection))  //安全遍历Conn
	Remove(IConnection)              //删除连接
	Clear()                          //删除并停止所有链接
	Count() int                      //获取当前连接
}
