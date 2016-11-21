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
	"sort"
	"strings"
	"fmt"
	"net"
	"github.com/FTwOoO/vpncore/common"

)

type IPList []net.IP

func (a IPList) Len() int {
	return len(a)
}
func (a IPList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a IPList) Less(i, j int) bool {
	return common.IPToUInt(a[i]) < common.IPToUInt(a[j])
}

func (a IPList) Sort() {
	sort.Sort(a)
}

func (self IPList) Contains(ip net.IP) bool {
	// TODO: support IPv6

	i := sort.Search(len(self), func(i int) bool {
		return common.IPToUInt(self[i]) >= common.IPToUInt(ip)
	})
	return i < len(self) && self[i].Equal(ip)
}

func (self *IPList) UnmarshalTOML(data []byte) (err error) {
	s := string(data)
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "[]")
	arr := strings.Split(s, ",")

	for _, ip := range arr {
		ip = strings.TrimSpace(ip)
		ip = strings.Trim(ip, "\"")
		if ip == "" {
			continue
		}

		ipobj := net.ParseIP(ip)
		if ipobj == nil {
			return fmt.Errorf("Bad IP format: %s", ip)
		} else {
			*self = append(*self, ipobj)
		}
	}

	self.Sort()
	return nil
}

