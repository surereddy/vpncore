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
	"github.com/golang/protobuf/proto"
	"io"
)

var (
	ErrInvalidArgs = errors.New("Invalid arguments.")
	ErrInvalidCtx = errors.New("Invalid context")
)

type Layer int

const (
	TRANSPORT_LAYER = Layer(1)
	OBS_LAYER = Layer(2)
	CRYPTO_LAYER = Layer(3)
	AUTH_LAYER = Layer(4)
	APPCATIOIN_LAYER = Layer(5)
)

type TransportProtocol string

const (
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
	NewListener(contexts []Context) (MessageListener, error)
}

type Client interface {
	//valid contexts pattern:
	//	StreamCreationContext -> StreamContext* -> MessageTransitionContext -> MessageContext*
	//	MessageCreationContext -> MessageContext*
	Dial(contexts []Context) (MessageIO, error)
}

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
	NewPipe(base StreamIO) StreamIO
}

type MessageCreationContext interface {
	Context
	Dial() (MessageIO, error)
	Listen() (MessageListener, error)
}

type MessageTransitionContext interface {
	Context
	NewPipe(base StreamIO) MessageIO
}

type MessageContext interface {
	Context
	NewPipe(base MessageIO) MessageIO
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
	io.ReadWriteCloser

	/*Read() (Message, error)
	Write(Message) error
	// Close closes the connection.
	// Any blocked Read or Write operations will be unblocked and return errors.
	Close() error

	// LocalAddr returns the local network address.
	LocalAddr() net.Addr

	// RemoteAddr returns the remote network address.
	RemoteAddr() net.Addr*/
}

type Message interface {
	proto.Marshaler
	proto.Unmarshaler
}

type Dialer interface {
	Dial() (net.Conn, error)
}

type DialerFunc func() (net.Conn, error)

func (f DialerFunc) Dial() (net.Conn, error) {
	return f()
}


/*

type Handler interface {
	HandleSession(*Session)
}

type HandlerFunc func(*Session)

func (f HandlerFunc) HandleSession(session *Session) {
	f(session)
}

func CreateCodec(dialer Dialer, protocol MessageContext) (MessagegReadWriteCloser, error) {
	conn, err := dialer.Dial()
	if err != nil {
		return nil, err
	}

	codec, err := protocol.NewCodec(conn)
	if err != nil {
		return nil, err
	}
	return codec, nil
}
*/
