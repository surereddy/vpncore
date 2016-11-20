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
)

type IPRange struct {
	Subnet *net.IPNet
	Start  uint32
	End    uint32
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

func (a IPNetList) Contains(ip net.IP) bool {
	ipval := common.IPToInt(ip.To4())

	l := len(a)
	i := sort.Search(l, func(i int) bool {
		n := a[i]
		return n.End >= ipval
	})

	if i < l {
		n := a[i]
		if n.Start <= ipval {
			return true
		}
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
			start := common.IPToInt(ipNet.IP)
			end := start + ^common.IPToInt(net.IP(ipNet.Mask))

			*self = append(*self, IPRange{Subnet:ipNet, Start:start, End:end})
		}
	}

	self.Sort()
	return nil
}










