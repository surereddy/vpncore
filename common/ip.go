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

package common

import "net"

import "encoding/binary"

func IPToUInt(ip net.IP) uint32 {
	// net.ParseIP will return 16 bytes for IPv4,
	// but we cant stop user creating 4 bytes for IPv4 bytes using net.IP{N,N,N,N}
	//
	if len(ip) == net.IPv4len {
		return binary.BigEndian.Uint32(ip)
	}

	return 0
}

func IP4FromUint32(v uint32) net.IP {
	ip := make([]byte, net.IPv4len)
	binary.BigEndian.PutUint32(ip, v)
	return ip
}