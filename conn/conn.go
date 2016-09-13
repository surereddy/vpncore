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

package conn

import (
	"net"
	"errors"
)

type TransProtocol uint32

const (
	PROTO_KCP = TransProtocol(1)
	PROTO_OBFS4 = TransProtocol(2)
)

type Listener struct {
	listener net.Listener
	proto TransProtocol
	//block cipher.Block
}

func NewListener(proto TransProtocol, addr net.Addr) (net.Listener, error) {
	switch proto {
	case PROTO_KCP:
		return &Listener{proto:proto, listener:&KCPListener{}}, nil
	default:
		return nil, errors.New("UNKOWN PROTOCOL!")
	}
}

func (l *Listener) Accept() (conn net.Conn, err error) {
	for {
		conn, err = l.listener.Accept()
		if err != nil {
			return
		}

	}
	return
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *Listener) Close() error {
	return l.listener.Close()
}

// Addr returns the listener's network address.
func (l *Listener) Addr() net.Addr {
	return l.listener.Addr()
}
