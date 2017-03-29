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

package transport

import (
	"net"
	"github.com/FTwOoO/vpncore/net/conn"
)

var _ conn.StreamListener = new(transportListener)

type transportListener struct {
	net.Listener
	proto conn.TransportProtocol
}

func (l *transportListener) Accept() (conn.StreamIO, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	} else {
		return &transportIO{Conn:c, proto:l.proto}, nil
	}
}