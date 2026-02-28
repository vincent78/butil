package pack

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const MaxPacketSize = 1024

type IPacket interface {
	GetHeadLen() int
	Pack(IMessage) ([]byte, error)
	UnPack([]byte) (IMessage, error)
}

type Packet struct{}

// NewPacket 封包拆包实例初始化方法
func NewPacket() IPacket {
	return &Packet{}
}

func (p *Packet) GetHeadLen() int {
	return 12
}

// Pack 打包方法(压缩数据)
func (p *Packet) Pack(msg IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	//写cmd
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetCmd()); err != nil {
		return nil, err
	}
	//写time
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetTime()); err != nil {
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
func (p *Packet) UnPack(binaryData []byte) (IMessage, error) {
	// 最大包长度
	if uint32(len(binaryData)) > MaxPacketSize {
		return nil, fmt.Errorf(`too large msg data received: %d`, len(binaryData))
	}

	//buffer
	dataBuff := bytes.NewReader(binaryData)

	//message
	msg := msgUpGetFromPool()

	//读cmd
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Cmd); err != nil {
		return nil, err
	}

	//读time
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Time); err != nil {
		return nil, err
	}

	//读data
	if cap(msg.Data) < len(binaryData)-p.GetHeadLen() {
		msg.Data = make([]byte, len(binaryData)-p.GetHeadLen())
	} else if cap(msg.Data) > len(binaryData)-p.GetHeadLen() {
		msg.Data = msg.Data[:len(binaryData)-p.GetHeadLen()]
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Data); err != nil {
		return nil, err
	}

	return msg, nil
}
