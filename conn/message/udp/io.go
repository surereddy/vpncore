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
	"sync"
)

type udpMessageIO struct {
	c          *net.UDPConn
	isServer   bool

	remoteAddr *net.UDPAddr

	ReadChan   chan []byte
	LastSeen   time.Time
	closeChan  chan struct{}
	closeOnce  sync.Once

}

func NewUdpMessageConn(c *net.UDPConn, isServer bool, remoteAddr *net.UDPAddr) (*udpMessageIO, error) {

	if (isServer && remoteAddr == nil) || c == nil {
		return nil, conn.ErrInvalidArgs
	}

	s := new(udpMessageIO)
	s.c = c
	s.isServer = isServer
	s.remoteAddr = remoteAddr
	s.ReadChan = make(chan []byte, 1024)
	s.closeChan = make(chan struct{})
	return s, nil
}

func (this *udpMessageIO) Read() ([]byte, error) {
	if this.isServer {
		buf := <-this.ReadChan
		return buf, nil
	} else {
		buf := make([]byte, 0xffff)
		n, _, err := this.c.ReadFromUDP(buf)
		if err != nil {
			return nil, err
		} else {
			return buf[:n], nil
		}
	}
}

func (this *udpMessageIO) Write(b []byte) (err error) {

	n := len(b)
	var wn int

	for {
		if !this.isServer {
			wn, err = this.c.Write(b)
		} else {
			wn, err = this.c.WriteToUDP(b, this.remoteAddr)
		}

		if err != nil {
			return
		}

		n -= wn
		if n > 0 {
			b = b[n:]
		} else {
			break
		}
	}

	return
}

func (this *udpMessageIO) Close() (err error) {
	this.closeOnce.Do(func() {

		close(this.closeChan)
		close(this.ReadChan)
		if !this.isServer {
			err = this.c.Close()
		}
	})

	return

}

func (this *udpMessageIO) LocalAddr() net.Addr {
	return this.c.LocalAddr()

}
func (this *udpMessageIO) RemoteAddr() net.Addr {
	if this.isServer {
		return this.remoteAddr
	} else {
		return this.RemoteAddr()
	}
}
