package ziface

type Packet interface {
	GetHeadLen() int
	Pack(IMessage) ([]byte, error)
	UnPack([]byte) (IMessage, error)
}
