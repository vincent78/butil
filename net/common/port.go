package common

import (
	"math/rand"
	"net"
	"strconv"
	"time"
)

// ScanPort 测试端口占用情况
func ScanPort(protocol string, hostname string, port int, timeout time.Duration) bool {
	p := strconv.Itoa(port)
	addr := net.JoinHostPort(hostname, p)
	conn, err := net.DialTimeout(protocol, addr, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func GetUnUsedPort(host string, num, min, max int) []int {
	tryNum := 0
	if max < 0 {
		max = 60000
	}

	if min < 0 {
		min = 50000
	}

	r := make([]int, 0, num)
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < num; i++ {
		v := rd.Intn(max-min) + min
		if !ScanPort("tcp", host, v, 1000*time.Millisecond) {
			r = append(r, v)
		} else {
			i = i - 1
			tryNum = tryNum + 1
			if tryNum >= 5 {
				return r
			}
		}
	}
	return r
}

// AllocatePorts allocates random ports suitable for listening.
func AllocatePorts(count int) ([]string, error) {
	var ports []string
	for i := 0; i < count; i++ {
		li, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, err
		}
		_, port, _ := net.SplitHostPort(li.Addr().String())
		err = li.Close()
		if err != nil {
			return nil, err
		}
		ports = append(ports, port)
	}
	return ports, nil
}

func GetFreePort() (port int, err error) {
	// 解析地址
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, nil
	}
	// 利用 ListenTCP 方法的如下特性
	// 如果 addr 的端口字段为0，函数将选择一个当前可用的端口
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, nil
	}
	// 关闭资源
	defer listen.Close()
	// 为了拿到具体的端口值，我们转换成 *net.TCPAddr类型获取其Port
	return listen.Addr().(*net.TCPAddr).Port, nil
}
