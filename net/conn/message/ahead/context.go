/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

package ahead

import (
	"github.com/FTwOoO/vpncore/net/conn"
	"crypto/aes"
	"crypto/cipher"
	"time"
	"encoding/binary"
	"crypto/rand"
	"errors"
)

func cipherAESGCM(k [32]byte) cipher.AEAD {
	c, err := aes.NewCipher(k[:])
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic(err)
	}
	return gcm
}

type AheadContext struct {
	cipher cipher.AEAD
	key    [32]byte
}

func NewAheadContext(key []byte) *AheadContext {
	ctx := new(AheadContext)
	copy(ctx.key[:], key)
	ctx.cipher = cipherAESGCM(ctx.key)
	return ctx
}

func (this *AheadContext) Valid() (bool, error) {
	return true, nil
}

func (this *AheadContext) Layer() conn.Layer {
	return conn.CRYPTO_LAYER
}

func (this *AheadContext) Encode(b []byte) ([]byte, error) {
	nonce := make([]byte, 12)
	n := time.Now().Unix()
	binary.BigEndian.PutUint64(nonce[4:], uint64(n))
	rand.Read(nonce[:4])

	en := this.cipher.Seal(nonce, nonce, b, nil)
	return en, nil
}

func (this *AheadContext) Decode(b []byte) ([]byte, error) {
	if len(b) < 12 {
		return nil, errors.New("Bad message to decode by Ahead")
	}

	nonce := b[:12]
	bts, err := this.cipher.Open(nil, nonce, b[12:], nil)
	return bts, err
}



