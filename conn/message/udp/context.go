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

package udp


import (
	"reflect"
	"net"
	"github.com/FTwOoO/vpncore/conn"
)

type UdpMessageContext struct {
}

func NewUdpMessageContext(uaddr string)  *UdpMessageContext {
		// Resolve UDP address
	uaddr2, err := net.ResolveUDPAddr("udp", uaddr)
	if err != nil {
		return nil, err
	}

	// Bind and setup UDP connection
	conn, err := net.ListenUDP("udp", uaddr2)
	if err != nil {
		return nil, err
	}

	// Start server
	s := &Server{
		conn: conn,
		peers: make(map[string]*client),
	}
	return &ProtobufMessageContext{msgTypes:msgTypes, rw:rw}
}


func (self *UdpMessageContext) Valid() (bool, error) {
	return true, nil
}


func (self *UdpMessageContext) Dial() (conn.MessageConn, error) {

}
func (self *UdpMessageContext) NewListener() (conn.MessageListener, error) {

}

