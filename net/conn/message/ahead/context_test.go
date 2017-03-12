/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

package ahead

import (
	"testing"
	"bytes"
)

func TestAheadContext(t *testing.T) {
	ctx:= NewAheadContext([]byte("Key..."))

	raw1 := []byte{1,2,3,4}
	en, err := ctx.Encode([]byte{1,2,3,4})
	if err != nil {
		t.Fatal("AheadContext encode error:", err)
	}

	raw2, err:= ctx.Decode(en)
	if err != nil {
		t.Fatal("AheadContext decode error:", err)
	}

	if !bytes.Equal(raw1, raw2) {
		t.Fatal("AheadContext encode/decode fail")
	}





}
