/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

package msgpack

import (
	"github.com/FTwOoO/vpncore/net/conn"
)

var _ conn.MessageToObjectContext = new(MsgpackContext)


//implements MessageToObjectContext
type MsgpackContext struct{}

func (this *MsgpackContext) Valid() (bool, error) {
	return true, nil
}

func (this *MsgpackContext) Layer() conn.Layer {
	return conn.ENCRYPTION_LAYER
}

func (this *MsgpackContext) Encode(obj interface{}) ([]byte, error) {
	if _, ok := obj.(Message); !ok {
		return nil, conn.ErrUnsupportType
	}

	v := wrapMessage{ContentMsg:obj.(Message)}
	return v.MarshalMsg(nil)
}

func (this *MsgpackContext) Decode(bts []byte) (obj interface{}, err error) {

	v := wrapMessage{}
	_, err = v.UnmarshalMsg(bts)
	if err != nil {
		return
	}
	return v.ContentMsg, nil
}



