/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

package msgpack

import (
	"github.com/tinylib/msgp/msgp"
	"testing"
	"bytes"
	"fmt"
	"encoding/hex"
)

func TestMarshalUnmarshalMsg(t *testing.T) {
	v := wrapMessage{ContentMsg:&TestMsg{Data:[]byte{1,2}}}
	bts, err := v.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(hex.Dump(bts))

	left, err := v.UnmarshalMsg(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after UnmarshalMsg(): %q", len(left), left)
	}

}

func TestEncodeDecodeMsg(t *testing.T) {
	v := wrapMessage{ContentMsg:&TestMsg{Data:[]byte{1,2}}}
	var buf bytes.Buffer
	msgp.Encode(&buf, &v)

	m := v.Msgsize()
	if buf.Len() > m {
		t.Logf("WARNING: Msgsize() for %v is inaccurate", v)
	}

	vn := wrapMessage{}
	err := msgp.Decode(&buf, &vn)
	if err != nil {
		t.Error(err)
	}
}
