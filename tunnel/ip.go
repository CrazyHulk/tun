package tunnel

import (
	"fmt"
	"net"
	"sync"
)

var maxIP uint16
var usedIPs = map[uint16]bool{}
var hostIPs = map[string]net.IP{}
var clientIPs = map[string]net.IP{}
var m sync.Mutex

func allocIPByTunName(tunName string) (host, client net.IP) {
	host = nextIP()
	client = nextIP()
	hostIPs[tunName] = host
	clientIPs[tunName] = client
	return
}

func releaseByTunName(name string) {
	hip, ok := hostIPs[name]
	if ok {
		releaseIP(hip)
	}
	cip, ok := clientIPs[name]
	if ok {
		releaseIP(cip)
	}
}

func releaseIP(ip net.IP) {
	m.Lock()
	defer m.Unlock()

	ip = ip.To4()

	i := (uint16(ip[2]) << 8) ^ uint16(ip[3])
	delete(usedIPs, i)
}

func nextIP() (ip net.IP) {
	m.Lock()
	defer m.Unlock()

	n := 0
	for ; maxIP == 0 || usedIPs[maxIP]; maxIP++ {
		n++
		if n&0x00ff == 0x00ff {
			continue
		}
		if n > 0xffff {
			maxIP = 0
			return
		}
	}

	usedIPs[maxIP] = true

	s := fmt.Sprintf("10.0.%d.%d", (maxIP>>8)&0xff, maxIP&0xff)

	ip = net.ParseIP(s)
	return
}
