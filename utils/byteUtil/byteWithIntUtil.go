package byteUtil

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

/***********************************************************************************************
主机字节序模式有两种，大端数据模式和小端数据模式
网络的数据是以大端数据模式进行交互，而我们的主机大多数以小端模式处理
两个主机在网络通信需要经过如下转换过程：主机字节序 —> 网络字节序 -> 主机字节序

大端模式：Big-Endian就是高位字节排放在内存的低地址端，低位字节排放在内存的高地址端
低地址 --------------------> 高地址
高位字节                     地位字节
小端模式：Little-Endian就是低位字节排放在内存的低地址端，高位字节排放在内存的高地址端
低地址 --------------------> 高地址
低位字节                     高位字节
什么是高位字节和低位字节
例如在32位系统中，257转换成二级制为：00000000 00000000 00000001 01100101，其中
00000001 | 01100101
高位字节     低位字节

字节顺序	描述	示例（十六进制 0x12345678）
大端序 (Big-Endian)	高位在前，低位在后	12 34 56 78
小端序 (Little-Endian)	低位在前，高位在后	78 56 34 12


x86/x64 架构（大多数 PC 和服务器）：通常使用小端序。
ARM 架构：通常是可配置的，但移动端多见小端序。
网络协议（TCP/IP）：为了让两台不同架构的机器能听懂对方在说什么，IP 协议族规定在传输层、网络层等协议头中，多字节整数必须按照大端序传输。

***********************************************************************************************/

func Int2Byte(n int) byte {
	if n > math.MaxUint8 {
		return 0
	} else {
		return Int2Bytes(n)[3]
	}
}

// int 转成 [2]byte  uint16
func Int2Byte2(curr int) []byte {
	bs := Int2Bytes(curr)
	return bs[2:]
}

// Int2Bytes 整形转换成字节
func Int2Bytes(n int) []byte {
	return Int2BytesByOrder(n, binary.BigEndian)
}

func Int2BytesByOrder(n int, order binary.ByteOrder) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, order, x)
	return bytesBuffer.Bytes()
}

// Bytes2Int 字节转换成整形
func Bytes2Int(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

func Bytes2IntByOrder(b []byte, order binary.ByteOrder) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, order, &x)
	return int(x)
}

// Int2BytesBigEndian int 转大端 []byte
func Int2BytesBigEndian(n int64, bytesLength byte) ([]byte, error) {
	switch bytesLength {
	case 1:
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, new(int8(n)))
		return bytesBuffer.Bytes(), nil
	case 2:
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, new(int16(n)))
		return bytesBuffer.Bytes(), nil
	case 3:
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, new(int32(n)))
		return bytesBuffer.Bytes()[1:], nil
	case 4:
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, new(int32(n)))
		return bytesBuffer.Bytes(), nil
	case 5:
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, new(n))
		return bytesBuffer.Bytes()[3:], nil
	case 6:
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, new(n))
		return bytesBuffer.Bytes()[2:], nil
	case 7:
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, new(n))
		return bytesBuffer.Bytes()[1:], nil
	case 8:
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, new(n))
		return bytesBuffer.Bytes(), nil
	}
	return nil, fmt.Errorf("Int2BytesBigEndian b param is invaild")
}

// Int2BytesLittleEndian int 转小端 []byte
func Int2BytesLittleEndian(n int64, bytesLength byte) ([]byte, error) {
	switch bytesLength {
	case 1:
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.LittleEndian, new(int8(n)))
		return bytesBuffer.Bytes(), nil
	case 2:
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.LittleEndian, new(int16(n)))
		return bytesBuffer.Bytes(), nil
	case 3:
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.LittleEndian, new(int32(n)))
		return bytesBuffer.Bytes()[0:3], nil
	case 4:
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.LittleEndian, new(int32(n)))
		return bytesBuffer.Bytes(), nil
	case 5:
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.LittleEndian, new(n))
		return bytesBuffer.Bytes()[0:5], nil
	case 6:
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.LittleEndian, new(n))
		return bytesBuffer.Bytes()[0:6], nil
	case 7:
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.LittleEndian, new(n))
		return bytesBuffer.Bytes()[0:7], nil
	case 8:
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.LittleEndian, new(n))
		return bytesBuffer.Bytes(), nil
	}
	return nil, fmt.Errorf("Int2BytesLittleEndian b param is invaild")
}

// Bytes2UIntBigEndian (大端) []byte 转 uint
func Bytes2UIntBigEndian(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0}, b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp uint8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp uint16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp uint32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "Bytes2Int bytes lenth is invaild!")
	}
}

// Bytes2IntBigEndian (大端) []byte 转 int
func Bytes2IntBigEndian(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0}, b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp int8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp int16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp int32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "Bytes2Int bytes lenth is invaild!")
	}
}

// Bytes2UIntLittleEndian (小端) []byte 转 uint
func Bytes2UIntLittleEndian(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0}, b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp uint8
		err := binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp uint16
		err := binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp uint32
		err := binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "Bytes2Int bytes lenth is invaild!")
	}
}

// Bytes2IntLittleEndian 小端[]byte 转 int
func Bytes2IntLittleEndian(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0}, b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp int8
		err := binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp int16
		err := binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp int32
		err := binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "Bytes2Int bytes lenth is invaild!")
	}
}

func Join(prefix, context []byte) []byte {
	result := make([]byte, len(prefix)+len(context))
	copy(result, prefix)
	copy(result[len(prefix):], context)
	return result
}

func AppendVarint(b []byte, v uint64) ([]byte, int) {
	s := 0
	for i := v; i > 0; s++ {
		i = i >> 8
	}
	switch {
	case v < 1<<7:
		b = append(b, byte(v))
	case v < 1<<14:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte(v>>7))
	case v < 1<<21:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte(v>>14))
	case v < 1<<28:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte(v>>21))
	case v < 1<<35:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte(v>>28))
	case v < 1<<42:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte(v>>35))
	case v < 1<<49:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte(v>>42))
	case v < 1<<56:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte(v>>49))
	case v < 1<<63:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte((v>>49)&0x7f|0x80),
			byte(v>>56))
	default:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte((v>>49)&0x7f|0x80),
			byte((v>>56)&0x7f|0x80),
			1)
	}
	return b, s
}

func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func uint64ToBytes(u uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, u)
	return buf
}
