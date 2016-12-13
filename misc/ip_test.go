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

package misc

import (
	"testing"
	"net"
)

func TestIP4FromAndFromUint32(t *testing.T) {

	ips := [] string{
		"1.2.3.4",
		"254.123.0.3",
		"127.0.0.1",
	}

	for _, ip := range ips {
		v := IP4ToUInt32(net.ParseIP(ip))
		ip2 := IP4FromUint32(v).String()

		if ip != ip2 {
			t.Failed()
		}

	}

}
