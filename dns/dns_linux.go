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
package vpncore

import (
	"net"
)

func (self *DNSManager) SetupNewDNS(new_dns []net.IP) (err error) {
	panic("Not implemented for this platform")

}

func (self *DNSManager) GetCurrentDNS() (l DNSList, err error) {

	panic("Not implemented for this platform")

}

func (self *DNSManager) RestoreDNS() (err error) {

	panic("Not implemented for this platform")
}