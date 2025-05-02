package byteUtil

import (
	"fmt"
	"github.com/imroc/biu"
)

/**
  1个字节=8个二进制位,每种数据类型占用的字节数都不一样
  注意位操作千万不要越界了，如某个类型占8个bit位，偏移时候不要超过这个范围
*/

// 最右边的 position = 0
func BitSet1(frame *byte, position int) {
	var data byte
	data = 1 << position
	*frame = (*frame) | data
}

func BitSet0(frame *byte, position int) {
	var data byte
	data = 1 << position
	*frame = (*frame) & (^data)
}

// BitGet 获取指定位的值  （是从最右边开始计算）
func BitGet(frame byte, position int) byte {
	return (frame << (7 - position)) >> 7
}

func Poll2Byte(frame *[]byte) int16 {
	length := len(*frame)
	if len(*frame) >= 2 {
		data := (*frame)[0:2]
		*frame = (*frame)[2:length]
		return (int16(data[0]) << 8) + int16(data[1])
	} else {
		return 0
	}

}
func Poll4Byte(frame *[]byte) int32 {
	length := len(*frame)
	if length >= 4 {
		data := (*frame)[0:4]
		*frame = (*frame)[4:length]
		return (int32(data[0]) << 24) + (int32(data[1]) << 16) + (int32(data[2]) << 8) + int32(data[3])
	} else {
		return 0
	}
}

// a输出结果:00011110  a = 30
func ShowBin(a byte) string {
	return biu.ToBinaryString(a)
}

/*
golang二进制bit位的常用操作，biu是一个转换二进制显示的库
*/
func test() {
	var (
		a uint8 = 30
	)
	//a输出结果:00011110
	fmt.Println(biu.ToBinaryString(a))
	/**
	  将某一位设置为1，例如设置第8位，从右向左数需要偏移7位,注意不要越界
	  1<<7=1000 0000 然后与a逻辑或|,偏移后的第8位为1，逻辑|运算时候只要1个为真就为真达到置1目的
	*/
	b := a | (1 << 7)
	//b输出结果:10011110
	fmt.Println(biu.ToBinaryString(b))
	/**
	  将某一位设置为0，例如设置第4位，从右向左数需要偏移3位,注意不要越界
	  1<<3=0000 1000 然后取反得到 1111 0111 然后逻辑&a
	*/
	c := a &^ (1 << 3)
	//c输出结果:00010110
	fmt.Println(biu.ToBinaryString(c))
	/**
	  获取某一位的值,即通过左右偏移来将将某位的值移动到第一位即可，当然也可以通过计算获得
	  如获取a的第4位
	  先拿掉4位以上的值 a<<4=1110 0000,然后拿掉右边的3位即可 a>>7=0000 0001
	*/
	d := (a << 4) >> 7
	//d输出结果:00000001 即1
	fmt.Println(biu.ToBinaryString(d))
	/**
	  取反某一位，即将某一位的1变0，0变1
	  这里使用到了亦或操作符 ^ 即 位值相同位0，不同为1
	  如获取a的第4位 1<<3=0000 1000
	  0000 1000 ^ 0001 1110 = 0001 0110
	*/
	e := a ^ (1 << 3)
	//d输出结果:00010110 即1
	fmt.Println(biu.ToBinaryString(e))

	/**
	  最后1个是综合用法,若tcp协议需要客户端先发送握手包，该包占用1个字节，其中前2位保留字段必须要为0，中间3位客户端对服务器版本要求，最后位客户端端版本
	  假设我们对服务器的版本要求和自己的版本都是3，那么我们该怎样构建这个包呢? 目标0001 1011
	  很多语言类型都没有直接 bit 单位，只要字节因此需要通过其他方法来得到,其实简单|或运算加上偏移即可,值得注意的网络使用的都是大端字节，传输前需要转换
	  rf=0 0000 0000
	  svf=3 0000 0011 偏移3位得到 0001 1000
	  cvf=3 0000 0011
	  计算
	  0000 0000
	  |
	  0001 1000
	  |
	  0000 0011
	  =
	  0001 1011
	*/
	var rf, svf, cvf uint8 = 0, 3, 3
	head := rf | (svf << 3) | cvf
	//head输出结果:00011011
	fmt.Println(biu.ToBinaryString(head))
}
