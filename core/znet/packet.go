package znet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/tonly18/xws/core/zconf"
	"github.com/tonly18/xws/core/ziface"
)

type Packet struct{}

// NewPacket 封包拆包实例初始化方法
func NewPacket() ziface.Packet {
	return &Packet{}
}

func (p *Packet) GetHeadLen() int {
	return 4
}

// Pack 打包方法(压缩数据)
func (p *Packet) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	//写cmd
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetCmd()); err != nil {
		return nil, err
	}
	//写data
	if data := msg.GetData(); len(data) > 0 {
		if err := binary.Write(dataBuff, binary.LittleEndian, data); err != nil {
			return nil, err
		}
	}

	return dataBuff.Bytes(), nil
}

// UnPack 拆包方法(解压数据)
func (p *Packet) UnPack(binaryData []byte) (ziface.IMessage, error) {
	// 最大包长度
	if zconf.Config.MaxPacketSize > 0 && uint32(len(binaryData)) > zconf.Config.MaxPacketSize {
		return nil, fmt.Errorf(`too large msg data received: %d`, len(binaryData))
	}

	//buffer
	dataBuff := bytes.NewReader(binaryData)

	//message
	msg := Message{}

	//读cmd
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Cmd); err != nil {
		return nil, err
	}

	//读data
	msg.Data = make([]byte, len(binaryData)-p.GetHeadLen())
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Data); err != nil {
		return nil, err
	}

	return &msg, nil
}
