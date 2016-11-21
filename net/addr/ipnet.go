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
}

func NewIPRangeByRange(start uint32, end uint32) *IPRange {
	//TODO: need ipv6 version?
	count := end - start
	if count <= 0 {
		return nil
	}

	ones := net.IPv4len * 8 - int(math.Floor(math.Log2(float64(count)) + 0.5))
	mark := net.CIDRMask(ones, net.IPv4len * 8)
	ip := common.IP4FromUint32(start)
	subnet := net.IPNet{IP: ip.Mask(mark), Mask: mark}
	return &IPRange{Start:start, End:end, Subnet:&subnet, version:IPv4}
}

func NewIPRangeByIPNet(subnet *net.IPNet) *IPRange {
	if ip4 := subnet.IP.To4(); ip4 != nil {
		start := common.IPToUInt(ip4)

		end := start + ^common.IPToUInt(ip4.Mask(subnet.Mask))
		return &IPRange{Subnet:subnet, Start:start, End:end, version:IPv4}
	}
	return nil
}

func (a *IPRange) Contains(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		if a.version == IPv4 {
			ipval := common.IPToUInt(ip4)
			return a.End >= ipval && a.Start <= ipval
		}
	}

	return false
}

type IPNetList []IPRange

func (a IPNetList) Len() int {
	return len(a)
}
func (a IPNetList) Less(i, j int) bool {
	return a[i].End < a[j].End
}
func (a IPNetList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a IPNetList) Sort() {
	sort.Sort(a)
}

// notice: used before Sort() called
func (a IPNetList) Contains(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {

		l := len(a)
		i := sort.Search(l, func(i int) bool {
			r := a[i]
			return r.Contains(ip4)
		})

		if i < l {
			return true
		}
		return false
	}

	return false
}

func (self *IPNetList) UnmarshalTOML(data []byte) error {
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
			r := *NewIPRangeByIPNet(ipNet)
			*self = append(*self, r)

		}
	}

	self.Sort()
	return nil
}










