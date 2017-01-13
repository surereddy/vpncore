package addr

import (
	"net"
	"strings"
	"fmt"
	"strconv"
)

type DNSAddresss struct {
	Ip   net.IP
	Port uint16
}

func (self *DNSAddresss) UnmarshalTOML(data []byte) (err error) {
	s := strings.Trim(string(data), "\"")
	spliter := ":"

	if strings.Contains(s, spliter) {
		arr := strings.Split(s, spliter)
		if len(arr) != 2 {
			return fmt.Errorf("Bad format for DNS:%s", s)
		}
		ip, port := arr[0], arr[1]
		self.Ip = net.ParseIP(ip)
		uport, err := strconv.Atoi(port)
		if err != nil {
			return err
		}

		self.Port = uint16(uport)
	} else {
		self.Ip = net.ParseIP(s)
		if self.Ip == nil {
			return fmt.Errorf("Bad format for DNS:%s", s)
		}

		self.Port = 53
	}
	return nil
}

func (self *DNSAddresss) String() string {
	return fmt.Sprintf("%s:%d", self.Ip.String(), self.Port)
}




