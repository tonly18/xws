package znet

// Message 消息(包结构:cmd data)
type Message struct {
	Cmd  uint32 `json:"cmd"`  //消息:cmd
	Data []byte `json:"data"` //消息:包体
}

// NewMessage 创建一个Message消息包
func NewMessage(cmd uint32, data []byte) *Message {
	return &Message{
		Cmd:  cmd,
		Data: data,
	}
}

// GetCmd 获取cmd数据
func (msg *Message) GetCmd() uint32 {
	return msg.Cmd
}

// SetCmd 设置cmd数据
func (msg *Message) SetCmd(cmd uint32) {
	msg.Cmd = cmd
}

// SetData 设置data数据
func (msg *Message) SetData(data []byte) {
	if len(data) > 0 {
		msg.Data = append(msg.Data, data...)
	}
}

// GetData 获取data数据
func (msg *Message) GetData() []byte {
	return msg.Data
}
