/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

package addr


import "net"

import (
	"encoding/binary"
	"bytes"
)

var v4InV6Prefix = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}


func IP4ToUInt32(ip net.IP) uint32 {
	// net.ParseIP will return 16 bytes for IPv4,
	// but we cant stop user creating 4 bytes for IPv4 bytes using net.IP{N,N,N,N}
	//
	if len(ip) == net.IPv4len {
		return binary.BigEndian.Uint32(ip)
	}

	if len(ip) == net.IPv6len && bytes.Equal(ip[:12], v4InV6Prefix) {
		return binary.BigEndian.Uint32(ip[12:])
	}

	return 0
}

func IP4FromUint32(v uint32) net.IP {
	ip := make([]byte, net.IPv4len)
	binary.BigEndian.PutUint32(ip, v)
	return ip
}
