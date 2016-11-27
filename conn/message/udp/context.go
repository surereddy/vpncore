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
	"net"
	"github.com/FTwOoO/vpncore/conn"
)

type UdpMessageContext struct {
	udpAddr *net.UDPAddr
}

func NewUdpMessageContext(uaddr string)  (*UdpMessageContext, error) {
	// Resolve UDP address
	uaddr2, err := net.ResolveUDPAddr("udp", uaddr)
	if err != nil {
		return nil, err
	}


	// Start server
	s := &UdpMessageContext{
		udpAddr: uaddr2,
	}
	return s, nil
}


func (self *UdpMessageContext) Valid() (bool, error) {
	return true, nil
}

func (self *UdpMessageContext) 	Layer() conn.Layer {
	return conn.TRANSPORT_LAYER
}


func (self *UdpMessageContext) Dial() (conn.MessagegReadWriteCloser, error) {
	srcAddr := &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 0}

	c, err := net.DialUDP("udp", srcAddr, self.udpAddr)
	if err != nil {
		return nil, err
	}

	return NewUdpMessageConn(c, nil, nil), nil
}

func (self *UdpMessageContext) Listener() (conn.MessageListener, error) {
	// Bind and setup UDP connection
	conn, err := net.ListenUDP("udp", uaddr2)
	if err != nil {
		return nil, err
	}
}

