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
	"time"
	"sync"
)

type UdpMessageListener struct {
	c               *net.UDPConn
	closed          chan struct{}
	connections     map[string]*udpMessageIO
	connectionsLock sync.Mutex
	newConnections  chan *udpMessageIO

	buf             []byte
}

func NewUdpMessageListener(udpAddr *net.UDPAddr, buf []byte) (l *UdpMessageListener, err error) {
	c, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	l = &UdpMessageListener{c:c,
		closed:make(chan struct{}),
		connections:map[string]*udpMessageIO{},
		connectionsLock:sync.Mutex{},
		newConnections: make(chan *udpMessageIO, ConnectionsChanSize),
	}
	if buf == nil {
		l.buf = make([]byte, MaxUDPPacketSize)
	} else {
		l.buf = buf
	}

	go l.listenLoop()
	go l.expireLoop()

	return
}

func (this *UdpMessageListener) expireLoop() {
	for {
		this.expire()
		time.Sleep(Lifetime)
	}
}

func (this *UdpMessageListener) expire() {
	this.connectionsLock.Lock()
	defer this.connectionsLock.Unlock()

	for id, cc := range this.connections {
		if cc.LastSeen.Add(Lifetime).After(time.Now()) {
			cc.Close()

			// Is it secure to delete entry inside the loop of map?
			this.connections[id] = nil
			continue
		}
	}
}

func (this *UdpMessageListener) listenLoop() {

	for {
		if this.closed == nil {
			return
		}

		n, addr, err := this.c.ReadFromUDP(this.buf[:])
		if err != nil || n >= MaxUDPPacketSize + 1 {
			return
		}

		msg := make([]byte, n)
		copy(msg, this.buf[:n])

		this.connectionsLock.Lock()
		if c, ok := this.connections[addr.String()]; !ok || c == nil {
			newConn, err := NewUdpMessageConn(
				this.c,
				true,
				addr,
			)

			if err != nil {
				this.connectionsLock.Unlock()
				continue

			}
			this.connections[addr.String()] = newConn
			//TODO: dont block if connections exceed the connection pool
			this.newConnections <- newConn
		}

		c := this.connections[addr.String()]
		c.ReadChan <- msg

		this.connectionsLock.Unlock()
	}
}

func (this *UdpMessageListener) Accept() (conn.MessageIO, error) {

	select {
	case <-this.closed:
		return nil, conn.ErrIOClosed
	case c := <-this.newConnections:
		return c, nil
	}
}

func (this *UdpMessageListener) Close() error {
	this.connectionsLock.Lock()
	defer this.connectionsLock.Unlock()

	for _, cc := range this.connections {
		cc.Close()
		continue

	}
	this.connections = nil
	close(this.closed)
	this.closed = nil
	return this.c.Close()
}

func (this *UdpMessageListener)  Addr() net.Addr {
	return this.c.LocalAddr()
}