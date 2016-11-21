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
	"testing"
	"net"
	"github.com/FTwOoO/vpncore/common"
	"fmt"
)

func TestIPRanges_Contains(t *testing.T) {
	ips := IPList{
		net.IP{32, 32, 32, 5},
		net.IP{17, 3, 2, 2},
		net.IP{22, 33, 44, 254},
		net.IP{22, 33, 44, 253},
	}

	netl := new(IPRanges)
	err := netl.UnmarshalTOML([]byte(`[
	    "32.32.32.0/24", "17.3.4.2/16", "22.33.44.253/30",
	    ]`))

	if err != nil {
		t.Fatal(err)
	}

	fmt.Print(netl)
	for _, ip := range ips {
		if netl.Contains(ip) != true {
			t.Fatalf("Ip(%d) %v is not in %v", ip, common.IPToUInt(ip), netl)
		}
	}

}


func TestNewIPRangeByRange(t *testing.T) {
	ip1:= net.IP{32, 32, 31, 5}
	ip2:= net.IP{32, 32, 31, 10}

	start := common.IPToUInt(ip1)
	end := common.IPToUInt(ip2)
	mid := (start + end) / 2
	ip3 := common.IP4FromUint32(mid)

	r := NewIPRangeByRange(start, end)
	fmt.Print(r.Subnet)
	if (!r.Subnet.Contains(ip1) || !r.Subnet.Contains(ip2) || !r.Subnet.Contains(ip3)||
		!r.Contains(ip1) || !r.Contains(ip2) || !r.Contains(ip3) ) {
		t.Failed()
	}

}

