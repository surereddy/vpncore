/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */


package msgpack

import (
	"github.com/tinylib/msgp/msgp"
	"reflect"
	"errors"
	"fmt"
	"sync"
)

var messages = make(map[MessageType]reflect.Type)
var messagesLock = sync.Mutex{}

func RegisterMessage(cmd MessageType, typ reflect.Type) {
	//insure reflect.New(typ).Interface().(MsgType)
	messagesLock.Lock()
	messages[cmd] = typ
	messagesLock.Unlock()
}

type MessageType uint8

const (
	MessageTypeTest = 0
	MessageTypeDNSRequest MessageType = 1
	MessageTypeDNSResponse MessageType = 2
	MessageTypeConnectionOpen MessageType = 3
	MessageTypeConnectionOpenDone MessageType = 4
	MessageTypeConnectionClose MessageType = 5
	MessageTypeConnectionData MessageType = 6
)

type Message interface {
	msgp.Sizer
	msgp.Decodable
	msgp.Encodable
	msgp.Marshaler
	msgp.Unmarshaler

	Cmd() MessageType
}

type wrapMessage struct {
	ContentMsg Message
}

// DecodeMsg implements msgp.Decodable
func (z *wrapMessage) DecodeMsg(dc *msgp.Reader) (err error) {
	var t uint8
	t, err = dc.ReadUint8()

	if err != nil {
		return
	}

	if typ, ok := messages[MessageType(t)]; ok {
		z.ContentMsg = reflect.New(typ).Interface().(Message)
		return z.ContentMsg.(msgp.Decodable).DecodeMsg(dc)
	} else {
		return errors.New("Invalid cmd")
	}
}

// EncodeMsg implements msgp.Encodable
func (z wrapMessage) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteUint8(uint8(z.ContentMsg.Cmd()))
	if err != nil {
		return
	}

	return z.ContentMsg.(msgp.Encodable).EncodeMsg(en)
}

// MarshalMsg implements msgp.Marshaler
func (z wrapMessage) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendUint8(o, uint8(z.ContentMsg.Cmd()))
	return z.ContentMsg.MarshalMsg(o)
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *wrapMessage) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zwht uint8
	zwht, bts, err = msgp.ReadUint8Bytes(bts)
	if err != nil {
		return
	}

	cmd := MessageType(zwht)
	fmt.Printf("cmd:%d", cmd)

	if typ, ok := messages[cmd]; ok {
		z.ContentMsg = reflect.New(typ).Interface().(Message)
		return z.ContentMsg.(msgp.Unmarshaler).UnmarshalMsg(bts)
	} else {
		return nil, errors.New("Invalid cmd")
	}
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z wrapMessage) Msgsize() (s int) {
	s = msgp.Uint8Size + z.ContentMsg.Msgsize()
	return
}
