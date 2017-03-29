/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

package encryption

import (
	"github.com/FTwOoO/vpncore/net/conn"
	"crypto/aes"
	"crypto/cipher"
	"time"
	"encoding/binary"
	"crypto/rand"
	"errors"
)

var _ conn.MessageContext = new(GCM256Context)


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

type GCM256Context struct {
	cipher cipher.AEAD
	key    [32]byte
}

func NewGCM256Context(key []byte) *GCM256Context {
	ctx := new(GCM256Context)
	copy(ctx.key[:], key)
	ctx.cipher = cipherAESGCM(ctx.key)
	return ctx
}

func (this *GCM256Context) Valid() (bool, error) {
	return true, nil
}

func (this *GCM256Context) Layer() conn.Layer {
	return conn.ENCRYPTION_LAYER
}

func (this *GCM256Context) Encode(b []byte) ([]byte, error) {
	nonce := make([]byte, 12)
	n := time.Now().Unix()
	binary.BigEndian.PutUint64(nonce[4:], uint64(n))
	rand.Read(nonce[:4])

	en := this.cipher.Seal(nonce, nonce, b, nil)
	return en, nil
}

func (this *GCM256Context) Decode(b []byte) ([]byte, error) {
	if len(b) < 12 {
		return nil, errors.New("Bad message to decode by Ahead")
	}

	nonce := b[:12]
	bts, err := this.cipher.Open(nil, nonce, b[12:], nil)
	return bts, err
}



