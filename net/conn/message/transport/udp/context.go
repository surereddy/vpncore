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
	"github.com/FTwOoO/vpncore/net/conn"
	"time"
)

const (
	MaxUDPPacketSize = 0x10000 - 1
	Lifetime = time.Duration(time.Second * 300)
	ConnectionsChanSize = 1024
	MessageChanSizePerConn = 1024
)

var _ conn.MessageCreationContext = new(UdpMessageContext)

type UdpMessageContext struct {
	udpAddr *net.UDPAddr
}

func NewUdpMessageContext(addr string) (*UdpMessageContext, error) {
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	s := &UdpMessageContext{
		udpAddr: uaddr,
	}
	return s, nil
}

func (this *UdpMessageContext) Valid() (bool, error) {
	return true, nil
}

func (this *UdpMessageContext) Layer() conn.Layer {
	return conn.TRANSPORT_LAYER
}

func (this *UdpMessageContext) Dial() (conn.MessageIO, error) {
	srcAddr := &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 0}
	//要验证一次read()不管读多少个字节都消掉一个UDP packet
	c, err := net.DialUDP("udp", srcAddr, this.udpAddr)
	if err != nil {
		return nil, err
	}

	return NewUdpMessageConn(c, false, nil)

}

func (this *UdpMessageContext) Listen() (conn.MessageListener, error) {
	return NewUdpMessageListener(this.udpAddr, nil)
}
