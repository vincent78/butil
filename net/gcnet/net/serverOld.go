package gcnet

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"github.com/vincent78/butil/logger/logger1"
	gcnet "github.com/vincent78/butil/net/gcnet/proto.pb"
	"net"
)

type NetHandler func(data interface{}, conn net.Conn)

var NetHandlerMap = make(map[string]NetHandler)

var ClientMap = make(map[net.Addr]net.Conn)

func RegistNetHandler(cmd string, f NetHandler) {
	NetHandlerMap[cmd] = f
}

func init() {
	RegistNetHandler("unknown", func(data interface{}, conn net.Conn) {
		logger1.Error("net[%v] receive unknown command. data: %v", conn.RemoteAddr(), data)
	})
	RegistNetHandler("stop", func(data interface{}, conn net.Conn) {
		logger1.Error("net[%v] receive unknown command. data: %v", conn.RemoteAddr(), data)
	})
}

// 接收消息
func readMessageOld(conn net.Conn) {
	logger1.Info("new connect: %v", conn.RemoteAddr())
	ClientMap[conn.RemoteAddr()] = conn
	reader := bufio.NewReader(conn)
	//读消息
	for {
		data, err := Decode(reader)
		if err == nil && len(data) > 0 {
			p := &gcnet.CmdRequest{}
			err := proto.Unmarshal(data, p)
			if err != nil {
				conn.Write([]byte("数据格式有误，请重新发"))
			}
			//buf= make([]byte, 10240,10240)//前两个字节表示本次请求body为长度
			logger1.Debug("receive %s %+v \n", conn.RemoteAddr(), p)
			cmd := p.Cmd
			if cmd == "" {
				cmd = "unknown"
			}
			if f, exist := NetHandlerMap[cmd]; exist {
				go f(p.Payload, conn)
			} else {
				NetHandlerMap["unknown"](p.Payload, conn)
			}
		}
		if err != nil {
			panic(err)
		}
	}
}

// Decode 解码数据包
func Decode(reader *bufio.Reader) ([]byte, error) {
	//读取数据包的长度（从包头获取）
	lenByte, err := reader.Peek(2) //读取前四个字节的数据
	if err != nil {
		return []byte{}, err
	}
	//转成Buffer对象,设置为从小端读取
	buff := bytes.NewBuffer(lenByte)

	var len uint16                                     //读取的数据大小，初始化为0
	err = binary.Read(buff, binary.LittleEndian, &len) //从小端读取
	if err != nil {
		return []byte{}, err
	}

	//读取消息
	pkg := make([]byte, int(len)+2)
	//Buffered返回缓冲区中现有的可读取的字节数
	if reader.Buffered() < int(len)+2 { //如果读取的包头的数据大小和读取到的不符合
		hr := 0
		for hr < int(len)+2 {
			l, err := reader.Read(pkg[hr:])
			if err != nil {
				return []byte{}, err
			}
			hr += l
		}
	} else {
		_, err := reader.Read(pkg)
		if err != nil {
			return []byte{}, err
		}
	}
	return pkg[2:], nil

}
