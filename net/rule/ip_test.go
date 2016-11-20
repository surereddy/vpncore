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

package rule

import (
	"testing"
	"net"
	"fmt"
	"math/rand"
	mtesting "github.com/FTwOoO/vpncore/testing"

)

func BenchmarkNewBlackIP(b *testing.B) {

	r := rand.New(rand.NewSource(0))

	ips := []net.IP{}
	for i := 0; i < b.N; i++ {
		ip := mtesting.RandomIPv4Address(r)
		ips = append(ips, ip)
	}

	iprule := NewIPBlocker(ips)

	for i := 0; i < b.N; i++ {
		ip := ips[i]
		f := iprule.FindIP(ip)
		fmt.Printf("Find ip %s: %v\n", ip, f)

		ip = mtesting.RandomIPv4Address(r)
		f = iprule.FindIP(ip)
		fmt.Printf("Find ip %s: %v\n", ip, f)
	}

}
