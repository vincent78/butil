package gcnet

import (
	"fmt"
	"testing"
	"time"
)

var (
	host = "172.31.236.40"
)

func TestPortUsed(t *testing.T) {
	r := ScanPort("tcp", host, 6378, 1*time.Second)
	fmt.Printf("port used : %v", r)
}

func TestRandomAvalidePorts(t *testing.T) {
	ips, err := AllocatePorts(10)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ips)
}

func TestRemoteRandomAvalidePorts(t *testing.T) {
	ps := GetUnUsedPort("192.168.0.122", 5, -1, -1)
	fmt.Println(ps)
}
