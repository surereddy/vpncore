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
	"github.com/FTwOoO/vpncore/conn"
	"net"
)

type UdpMessageConn struct {
	udpConn    *net.UDPConn
	remoteAddr *net.UDPAddr
	buf        []byte
}

func (this *UdpMessageConn) Receive() (conn.Message, error) {

	this.buf = this.buf[:cap(this.buf)]
	n, _, err := this.udpConn.ReadFromUDP(this.buf)
	if err != nil {
		return nil, err
	}

	this.buf = this.buf[:n]

	msg := conn.Message{}
	msg.Unmarshal(this.buf)

}
func (this *UdpMessageConn) Send(msg conn.Message) (err error) {
	buf, err := msg.Marshal()
	if err != nil {
		return
	}

	if this.remoteAddr == nil {
		this.udpConn.Write(buf)
	} else {
		this.udpConn.WriteToUDP(buf, this.remoteAddr)
	}

}
func (this *UdpMessageConn) Close() error {
	this.udpConn.Close()

}
func (this *UdpMessageConn) LocalAddr() net.Addr {
	return this.udpConn.LocalAddr()

}
func (this *UdpMessageConn) RemoteAddr() net.Addr {
	addr1 := this.udpConn.RemoteAddr()
	if addr1 == "" {
		return this.remoteAddr
	}
}
