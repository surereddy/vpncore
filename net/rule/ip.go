package rule

import (
	"sort"
	"net"
	"github.com/FTwOoO/vpncore/net/addr"

)

type IPBlocker struct {
	Ip  addr.IPList
	Net addr.IPNetList
}

func NewIPBlocker(ips addr.IPList) (*IPBlocker) {
	sort.Sort(ips)
	return &IPBlocker{Ip:ips}
}

func (self *IPBlocker) FindIP(ip net.IP) bool {
	return self.Ip.Contains(ip) || self.Net.Contains(ip)
}