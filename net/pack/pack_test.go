package pack

import (
	"fmt"
	"testing"
)

func TestPacket_Pack(t *testing.T) {

	msg := NewMessage(1234567890, []byte("hello world"))

	pk := NewPacket()
	msgPacket, err := pk.Pack(msg)

	fmt.Println("<-- err:", err)
	fmt.Println("<-- msgPacket:", msgPacket)
	fmt.Println("<-- length:", len(msgPacket))
	fmt.Println("<-- msgPacket:", string(msgPacket))

	msg2, err := pk.UnPack(msgPacket)
	fmt.Println("--> err:", err)
	fmt.Println("--> cmd:", msg2.GetCmd())
	fmt.Println("--> time:", msg2.GetTime())
	fmt.Println("--> length:", len(msg2.GetData()))
	fmt.Println("--> data:", string(msg2.GetData()))
}
