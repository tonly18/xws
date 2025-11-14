package ziface

type IConnManager interface {
	Add(IConnection)                 //添加链接
	Get(uint64) (IConnection, error) //利用ConnID获取链接
	Remove(IConnection)              //删除连接
	Clear()                          //删除并停止所有链接

	GetByUid(int64) (IConnection, error)
	Len() (int, int) //获取当前连接
}
