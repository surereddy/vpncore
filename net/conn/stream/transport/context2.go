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

var _ conn.StreamCreationContext = new(TransportStreamContext2)

type TransportStreamContext2 struct {
	Listener net.Listener
}

func (this *TransportStreamContext2) Dial() (conn.StreamIO, error) {
	return nil, errors.New("ClientMode not supported!")
}

func (this *TransportStreamContext2) Listen() (conn.StreamListener, error) {
	return &transportListener{proto:conn.PROTO_UNKOWN, Listener:this.Listener}, nil
}

func (this *TransportStreamContext2) Layer() conn.Layer {
	return conn.TRANSPORT_LAYER
}

func (this *TransportStreamContext2) Valid() (bool, error) {
	return true, nil
}
