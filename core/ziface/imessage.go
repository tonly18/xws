package ziface

type IMessage interface {
	GetCmd() uint32 //获取cmd数据
	SetCmd(uint32)  //设置cmd数据

	SetData([]byte)  //设置data数据
	GetData() []byte //获取data数据
}
