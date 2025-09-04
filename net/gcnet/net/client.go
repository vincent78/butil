package gcnet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/vincent78/butil/net/gcnet/proto.pb"
)

func Client() {
	strIP := "localhost:6600"
	var conn net.Conn
	var err error

	//连接服务器
	if conn, err = net.Dial("tcp", strIP); err != nil {
		panic(err)
	}

	fmt.Println("connect", strIP, "success")
	defer conn.Close()

	//发送消息
	//sender := bufio.NewScanner(os.Stdin)
	for i := 0; i < 10000; i++ {
		stSend := &gcnet.CmdRequest{
			Cmd:     "echo",
			Payload: "send: " + strconv.Itoa(i),
		}
		//protobuf编码
		pData, err := proto.Marshal(stSend)
		if err != nil {
			panic(err)
		}
		data, err := Encode(pData)
		if err == nil {
			conn.Write(data)
			fmt.Printf("%d send %d \n", i, uint16(len(pData)))
		}

		/*if sender.Text() == "stop" {
		   return
		}*/
	}
	for {
		time.Sleep(time.Millisecond * 10)
		/*stSend := &model.People{
		     Config: "lishang "+ strconv.Itoa(22222222),
		     Age:  *proto.Int(22222222),
		  }

		  //protobuf编码
		  pData, err := proto.Marshal(stSend)
		  if err != nil {
		     panic(err)
		  }
		  b := make([]byte,2,2)
		  binary.BigEndian.PutUint16(b,uint16(len(pData)))
		  //发送
		  conn.Write(append(b,pData...))*/
	}

}

// Encode 将数据包编码（即加上包头再转为二进制）
func Encode(mes []byte) ([]byte, error) {
	//获取发送数据的长度，并转换为四个字节的长度，即int32
	len := uint16(len(mes))
	//创建数据包
	dataPackage := new(bytes.Buffer) //使用字节缓冲区，一步步写入性能更高

	//先向缓冲区写入包头
	//大小端口诀：大端：尾端在高位，小端：尾端在低位
	//编码用小端写入，解码也要从小端读取，要保持一致
	err := binary.Write(dataPackage, binary.LittleEndian, len) //往存储空间小端写入数据
	if err != nil {
		return nil, err
	}
	//写入消息
	err = binary.Write(dataPackage, binary.LittleEndian, mes)
	if err != nil {
		return nil, err
	}
	return dataPackage.Bytes(), nil
}
