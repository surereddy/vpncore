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
	"time"
	"github.com/FTwOoO/vpncore/conn"
)

type udpMessageIO struct {
	localAddr  *net.UDPAddr
	remoteAddr *net.UDPAddr
	ReadChan   chan []byte
	WriteChan  chan []byte
	LastSeen   time.Time
	Closed     chan struct{}
}

func NewUdpMessageConn(localAddr *net.UDPAddr, remoteAddr *net.UDPAddr, readChan <- chan []byte, writeChan<- chan []byte) (*udpMessageIO, error) {

	c := new(udpMessageIO)
	c.remoteAddr = remoteAddr
	c.localAddr = localAddr
	c.ReadChan = readChan
	c.WriteChan = writeChan
	c.Closed = make(chan struct{})
	return c, nil
}


func (this *udpMessageIO) Read(msg conn.Message) error {
	buf := <-this.ReadChan
	_, err := msg.Unmarshal(buf)
	return err

}

func (this *udpMessageIO) Write(msg conn.Message) (err error){
	b, err := msg.Marshal()
	if err != nil {
		return err
	}
	this.WriteChan <- b
	return len(b), nil
}

func (this *udpMessageIO) Close() error {
	close(this.Closed)
}

func (this *udpMessageIO) LocalAddr() net.Addr {
	return this.localAddr

}
func (this *udpMessageIO) RemoteAddr() net.Addr {
	return this.remoteAddr
}
