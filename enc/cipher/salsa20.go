package cipher

import (
	"golang.org/x/crypto/salsa20"
)


// Salsa20BlockCrypt implements BlockCrypt
type Salsa20Stream struct {
	key [32]byte
}

// NewSalsa20BlockCrypt initates BlockCrypt by the given key
func NewSalsa20Stream(key []byte) (StreamCipher, error) {
	c := new(Salsa20Stream)
	copy(c.key[:], key)
	return c, nil
}

// Encrypt implements Encrypt interface
func (c *Salsa20Stream) Encrypt(dst, src []byte) {
	salsa20.XORKeyStream(dst[:], src[:], c.key[:8], &c.key)
	copy(dst, src)
}

// Decrypt implements Decrypt interface
func (c *Salsa20Stream) Decrypt(dst, src []byte) {
	salsa20.XORKeyStream(dst[:], src[:], c.key[:8], &c.key)
	copy(dst, src)
}
