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
	"errors"
	"net"
	"strings"
	"fmt"
	"io"
)

var (
	ErrInvalidArgs = errors.New("Invalid arguments.")
	ErrInvalidCtx = errors.New("Invalid context")
	ErrIOClosed = errors.New("IO closed")
	ErrDecode = errors.New("Decode error")
	ErrInValidHandshakeStep = errors.New("Invalid handshake step")
	ErrUnsupportType = errors.New("Unsupported type")
)

type Layer int

const (
	TRANSPORT_LAYER = Layer(1)
	COMPRESS_LAYER = Layer(2)
	OBS_LAYER = Layer(3)
	ENCRYPTION_LAYER = Layer(4)
	FRAGMENT_LAYER = Layer(5)
	APPCATIOIN_LAYER = Layer(6)
)

type TransportProtocol string

const (
	PROTO_UNKOWN = TransportProtocol("unkown")
	PROTO_TCP = TransportProtocol("tcp")
	PROTO_KCP = TransportProtocol("kcp")
	PROTO_HTTP2 = TransportProtocol("http2")
	PROTO_QUIC = TransportProtocol("quic")
)

func (self *TransportProtocol) UnmarshalTOML(data []byte) (err error) {
	name := string(data)
	name = strings.TrimSpace(name)
	name = strings.Trim(name, "\"")

	switch name {
	case string(PROTO_TCP): *self = PROTO_TCP
	case string(PROTO_KCP): *self = PROTO_KCP
	case string(PROTO_HTTP2): *self = PROTO_HTTP2
	case string(PROTO_QUIC): *self = PROTO_QUIC

	default:
		return fmt.Errorf("invalid protocal:%s", name)
	}
	return
}

type Server interface {
	NewListener(contexts []Context) (ObjectListener, error)
}

type Client interface {
	//valid contexts pattern:
	//	StreamCreationContext -> StreamContext* -> MessageTransitionContext -> MessageContext*
	//	MessageCreationContext -> MessageContext*
	Dial(contexts []Context) (ObjectIO, error)
}

//A server or client own context instances. if a context object have states,
//the states is shared by all io connections, for example, NoiseIKMessageContext
//has the handshake state that changes when processing messages, no matter which
// connection the messages arride.
//
type Context interface {
	Valid() (bool, error)
	Layer() Layer
}

type StreamCreationContext interface {
	Context
	Dial() (StreamIO, error)
	Listen() (StreamListener, error)
}

type StreamContext interface {
	Context
	Pipe(base StreamIO) StreamIO
}

type StreamToMessageContext interface {
	Context
	Pipe(base StreamIO) MessageIO
}

type MessageCreationContext interface {
	Context
	Dial() (MessageIO, error)
	Listen() (MessageListener, error)
}

type MessageContext interface {
	Context
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}

type MessageToObjectContext interface {
	Context
	Encode(interface{}) ([]byte, error)
	Decode([]byte) (interface{}, error)
}

type  StreamListener interface {
	Accept() (StreamIO, error)
	Close() error
	Addr() net.Addr
}

type  StreamIO interface {
	io.ReadWriteCloser
}

type  MessageListener interface {
	Accept() (MessageIO, error)
	Close() error
	Addr() net.Addr
}

type MessageIO interface {
	Read() ([]byte, error)
	Write([]byte) error
	io.Closer
}

type  ObjectListener interface {
	Accept() (ObjectIO, error)
	Close() error
	Addr() net.Addr
}

type ObjectIO interface {
	Read() (interface{}, error)
	Write(interface{}) error
	io.Closer
}

type Marshaler interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) (int, error)
}

type Dialer interface {
	Dial() (net.Conn, error)
}

type DialerFunc func() (net.Conn, error)

func (f DialerFunc) Dial() (net.Conn, error) {
	return f()
}

