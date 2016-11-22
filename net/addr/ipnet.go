/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * Author: FTwOoO <booobooob@gmail.com>
 */

package addr

import (
	"net"
	"strings"
	"fmt"
	"sort"
	"github.com/FTwOoO/vpncore/common"
	"math"
)

type AddressType uint

const (
	IPv4 = AddressType(4)
	IPv6 = AddressType(6)
	Domain = AddressType(7)
)

type IPRange struct {
	version AddressType
	Subnet  *net.IPNet
	Start   uint32
	End     uint32
	Info    interface{}
}

func NewIPRangeByRange(start uint32, end uint32) *IPRange {
	//TODO: need ipv6 version?
	count := end - start + 1
	if count <= 0 {
		return nil
	}

	ones := net.IPv4len * 8 - int(math.Floor(math.Log2(float64(count)) + 0.5))
	mask := net.CIDRMask(ones, net.IPv4len * 8)
	ip := common.IP4FromUint32(start)
	subnet := net.IPNet{IP: ip.Mask(mask), Mask: mask}
	//fmt.Printf("start %x, ip %s subnet %s mask %d\n", start, ip, subnet, mask)
	return &IPRange{Start:start, End:end, Subnet:&subnet, version:IPv4}
}

func NewIPRangeByStartIp(ip net.IP, count uint32) *IPRange {
	if ip4 := ip.To4(); ip4 != nil {
		start := common.IP4ToUInt32(ip4)
		end := start + count - 1
		return NewIPRangeByRange(start, end)
	}
	return nil
}

func NewIPRangeByIPNet(subnet *net.IPNet) *IPRange {
	if ip4 := subnet.IP.To4(); ip4 != nil {
		maskedIp := subnet.IP.Mask(subnet.Mask)
		start := common.IP4ToUInt32(maskedIp)
		end := start + ^common.IP4ToUInt32(net.IP(subnet.Mask))

		if end < start {
			return nil
		}
		return &IPRange{Subnet:subnet, Start:start, End:end, version:IPv4}
	}
	return nil
}

func (a *IPRange) UpdateInfo(info interface{}) *IPRange {
	a.Info = info
	return a
}

func (a *IPRange) Contains(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		if a.version == IPv4 {
			ipval := common.IP4ToUInt32(ip4)
			return a.End >= ipval && a.Start <= ipval
		}
	}

	return false
}

func (a *IPRange) Less(b *IPRange) bool {
	return a.End < b.Start
}

type IPRanges []*IPRange

func (a IPRanges) Len() int {
	return len(a)
}
func (a IPRanges) Less(i, j int) bool {
	return a[i].Less(a[j])
}
func (a IPRanges) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a IPRanges) Sort() {
	sort.Sort(a)
}

// notice: used before Sort() called
func (a IPRanges) Contains(ip net.IP) bool {

	n := len(a)
	i := a.search(ip)

	if i < n {
		return true
	}

	return false
}

func (a IPRanges) Get(ip net.IP) *IPRange {

	n := len(a)
	i := a.search(ip)

	if i < n {
		return a[i]
	}
	return nil

}

func (a IPRanges) search(ip net.IP) int {
	if ip4 := ip.To4(); ip4 != nil {
		ipval := common.IP4ToUInt32(ip4)
		fmt.Printf("Start search %d \n", ipval)
		n := len(a)

		i, j := 0, n
		for i < j {
			h := i + (j - i) / 2 // avoid overflow when computing h
			r := a[h]

			if ipval > r.End {
				i = h + 1
			} else if ipval < r.Start {
				j = h
			} else {
				return h
			}
		}
	}

	return -1
}

func (self *IPRanges) UnmarshalTOML(data []byte) error {
	s := string(data)
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "[]")
	arr := strings.Split(s, ",")

	for _, subnet := range arr {
		subnet = strings.TrimSpace(subnet)
		subnet = strings.Trim(subnet, "\"")
		if subnet == "" {
			continue
		}

		_, ipNet, err := net.ParseCIDR(subnet)
		if err != nil {
			fmt.Println("ERR!")
			return err
		} else {
			r := NewIPRangeByIPNet(ipNet)
			*self = append(*self, r)

		}
	}

	self.Sort()
	return nil
}










