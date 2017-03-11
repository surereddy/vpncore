/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */


package crypto



import (
	"github.com/FTwOoO/noise"
	"testing"
	"bytes"
)


func TestAhead(t *testing.T) {

	key := [32]byte{1,3,5,7,9,100}
	plaintext := []byte{11,33,55,77,99}

	cf := noise.CipherAESGCM.Cipher(key)


	out := cf.Encrypt(nil, 1, []byte{33}, plaintext)

	out2, err := cf.Decrypt(nil, 1, []byte{33}, out)

	if err != nil {
		t.Fatal("decrypt fail!", err)
	}

	if !bytes.Equal(plaintext, out2) {
		t.Fatal("after encrypt end decrypt, plaintext will be same")
	}
}