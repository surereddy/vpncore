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
	"errors"
	"github.com/FTwOoO/vpncore/net/conn"
)

var _ conn.StreamCreationContext = new(TCPTransportStreamContext)


type TCPTransportStreamContext struct {
	Protocol   conn.TransportProtocol
	ListenAddr string
	RemoveAddr string
}

func (this *TCPTransportStreamContext) Dial() (conn.StreamIO, error){

	switch this.Protocol {
	case conn.PROTO_TCP:
		c, err := net.Dial("tcp", this.RemoveAddr)
		if err != nil {
			return nil, err
		}

		return &transportIO{Conn:c, proto:this.Protocol}, nil


	case conn.PROTO_KCP:
		panic("not implemented!")
	}

	return nil, errors.New("Proto not supported!")
}

func (this *TCPTransportStreamContext) Listen() (conn.StreamListener, error) {
	switch this.Protocol {
	case conn.PROTO_KCP:
		panic("not implemented yet!")
	case conn.PROTO_TCP:
		addr, err := net.ResolveTCPAddr("tcp4", this.ListenAddr)
		if err != nil {
			return nil, err
		}

		l, err := net.ListenTCP("tcp4", addr)
		if err != nil {
			return nil, err
		}
		return &transportListener{proto:this.Protocol, Listener:l}, nil
	default:
		return nil, errors.New("UNKOWN PROTOCOL!")
	}
}

func (this *TCPTransportStreamContext) Layer() conn.Layer {
	return conn.TRANSPORT_LAYER
}

func (this *TCPTransportStreamContext) Valid() (bool, error) {
	res := this.Protocol != "" && this.ListenAddr != "" && this.RemoveAddr != ""

	return res, nil
}
