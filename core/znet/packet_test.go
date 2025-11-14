package znet

import (
	"fmt"
	"testing"
)

func TestDataPack(t *testing.T) {
	msg := NewMessage(1234567890, []byte("hello world"))

	pk := NewPacket()
	msgPacket, err := pk.Pack(msg)

	fmt.Println("err:", err)
	fmt.Println("msgPacket:", msgPacket)
	fmt.Println("msgPacket:", string(msgPacket), len(msgPacket))

	msg2, err := pk.UnPack(msgPacket)
	fmt.Println("err:", err)
	fmt.Println("msg2:", msg2.GetCmd())
	fmt.Println("msg2:", string(msg2.GetData()), len(msg2.GetData()))
}
