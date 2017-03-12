/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

//go:generate msgp

package msgpack

import (
	"reflect"
)

func init() {
	RegisterMessage(MessageTypeTest, reflect.TypeOf(TestMsg{}))
}


//msgp:tuple TestMsg
type TestMsg struct {
	Data []byte `msg:"data"`
}

func (z TestMsg) Cmd() MessageType {
	return MessageTypeTest
}