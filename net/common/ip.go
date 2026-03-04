package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"strings"
)

func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := GetIp4FromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network")
}

func GetIp4FromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func GetIp6FromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To16()
	if ip == nil {
		return nil
	}
	return ip
}

// 获取外网ip地址
func GetLocation(ip, key string) string {
	if ip == "127.0.0.1" || ip == "localhost" {
		return "内部IP"
	}
	url := "https://restapi.amap.com/v5/ip?ip=" + ip + "&type=4&key=" + key
	fmt.Println("url", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("restapi.amap.com failed:", err)
		return "未知位置"
	}
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(s))

	m := make(map[string]string)

	err = json.Unmarshal(s, &m)
	if err != nil {
		fmt.Println("Umarshal failed:", err)
	}
	//if m["province"] == "" {
	//	return "未知位置"
	//}
	return m["country"] + "-" + m["province"] + "-" + m["city"] + "-" + m["district"] + "-" + m["isp"]
}

// 获取局域网ip地址
func GetLocaHonst() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}

	}
	return ""
}

// 获取容器内 IP
func LocalIP() string {
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok &&
				!ipnet.IP.IsLoopback() &&
				ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func ParseInterfaceAddr(ifceName, network string) (ifce string, addr []net.Addr, err error) {
	if ifceName == "" {
		addr = append(addr, nil)
		return
	}

	ip := net.ParseIP(ifceName)
	if ip == nil {
		var ife *net.Interface
		ife, err = net.InterfaceByName(ifceName)
		if err != nil {
			return
		}
		var addrs []net.Addr
		addrs, err = ife.Addrs()
		if err != nil {
			return
		}
		if len(addrs) == 0 {
			err = fmt.Errorf("addr not found for interface %s", ifceName)
			return
		}
		ifce = ifceName
		for _, addr_ := range addrs {
			if ipNet, ok := addr_.(*net.IPNet); ok {
				addr = append(addr, ipToAddr(ipNet.IP, network))
			}
		}
	} else {
		ifce, err = findInterfaceByIP(ip)
		if err != nil {
			return
		}
		addr = []net.Addr{ipToAddr(ip, network)}
	}

	return
}

func ipToAddr(ip net.IP, network string) (addr net.Addr) {
	port := 0
	switch network {
	case "tcp", "tcp4", "tcp6":
		addr = &net.TCPAddr{IP: ip, Port: port}
		return
	case "udp", "udp4", "udp6":
		addr = &net.UDPAddr{IP: ip, Port: port}
		return
	default:
		addr = &net.IPAddr{IP: ip}
		return
	}
}

func findInterfaceByIP(ip net.IP) (string, error) {
	ifces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, ifce := range ifces {
		addrs, _ := ifce.Addrs()
		if len(addrs) == 0 {
			continue
		}
		for _, addr := range addrs {
			ipAddr, _ := addr.(*net.IPNet)
			if ipAddr == nil {
				continue
			}
			// logger.Default().Infof("%s-%s", ipAddr, ip)
			if ipAddr.IP.Equal(ip) {
				return ifce.Name, nil
			}
		}
	}
	return "", nil
}

// AddrPortRange is the network address with port range supported.
// e.g. 192.168.1.1:0-65535
type AddrPortRange string

func (p AddrPortRange) Addrs() (addrs []string) {
	// ignore url scheme, e.g. http://, tls://, tcp://.
	if strings.Contains(string(p), "://") {
		return nil
	}

	h, sp, err := net.SplitHostPort(string(p))
	if err != nil {
		return nil
	}

	pr := PortRange{}
	pr.Parse(sp)

	for i := pr.Min; i <= pr.Max; i++ {
		addrs = append(addrs, net.JoinHostPort(h, strconv.Itoa(i)))
	}
	return addrs
}

// Port range is a range of port list.
type PortRange struct {
	Min int
	Max int
}

// Parse parses the s to PortRange.
// The s can be a single port number and will be converted to port range port-port.
func (pr *PortRange) Parse(s string) error {
	minmax := strings.Split(s, "-")
	switch len(minmax) {
	case 1:
		port, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		if port < 0 || port > 65535 {
			return fmt.Errorf("invalid port: %s", s)
		}

		pr.Min, pr.Max = port, port
		return nil

	case 2:
		min, err := strconv.Atoi(minmax[0])
		if err != nil {
			return err
		}
		max, err := strconv.Atoi(minmax[1])
		if err != nil {
			return err
		}

		pr.Min, pr.Max = min, max
		return nil

	default:
		return fmt.Errorf("invalid port range: %s", s)
	}
}

func (pr *PortRange) Contains(port int) bool {
	return port >= pr.Min && port <= pr.Max
}

type IPRange struct {
	Min netip.Addr
	Max netip.Addr
}

func (r *IPRange) Parse(s string) error {
	minmax := strings.Split(s, "-")
	switch len(minmax) {
	case 1:
		addr, err := netip.ParseAddr(strings.TrimSpace(s))
		if err != nil {
			return err
		}

		r.Min, r.Max = addr, addr
		return nil

	case 2:
		min, err := netip.ParseAddr(strings.TrimSpace(minmax[0]))
		if err != nil {
			return err
		}
		max, err := netip.ParseAddr(strings.TrimSpace(minmax[1]))
		if err != nil {
			return err
		}

		r.Min, r.Max = min, max
		return nil

	default:
		return fmt.Errorf("invalid ip range: %s", s)
	}
}

func (r *IPRange) Contains(addr netip.Addr) bool {
	return !(addr.Less(r.Min) || r.Max.Less(addr))
}
