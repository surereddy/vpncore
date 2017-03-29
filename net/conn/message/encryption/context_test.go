/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

package encryption

import (
	"testing"
	"bytes"
)

func TestGCM256Context(t *testing.T) {
	ctx := NewGCM256Context([]byte("Key..."))

	raw1 := []byte{1, 2, 3, 4}
	en, err := ctx.Encode([]byte{1, 2, 3, 4})
	if err != nil {
		t.Fatal("AheadContext encode error:", err)
	}

	raw2, err := ctx.Decode(en)
	if err != nil {
		t.Fatal("AheadContext decode error:", err)
	}

	if !bytes.Equal(raw1, raw2) {
		t.Fatal("AheadContext encode/decode fail")
	}

}
