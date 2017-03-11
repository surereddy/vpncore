/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

package msgpack

import (
	"github.com/FTwOoO/vpncore/net/conn"
)

type MsgpackContext struct {
}

func NewMsgpackContext(key []byte) *MsgpackContext {
	ctx := new(MsgpackContext)
	return ctx
}

func (this *MsgpackContext) Valid() (bool, error) {
	return true, nil
}

func (this *MsgpackContext) Layer() conn.Layer {
	return conn.CRYPTO_LAYER
}

func (this *MsgpackContext) Encode(b []byte) ([]byte, error) {
	return nil, nil
}

func (this *MsgpackContext) Decode(b []byte) ([]byte, error) {
	return nil, nil
}
