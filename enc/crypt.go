package enc

import (
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha1"
	"errors"


	//TODO: More elegant style to import these types from sub modules to
	// local namespace?
	"github.com/FTwOoO/vpncore/enc/cipher"
)


//TODO: More elegant style to import these types from sub modules to
// local namespace?
type CommonCipher cipher.CommonCipher
type StreamCipher cipher.StreamCipher
type BlockCipher cipher.BlockCipher

type Cipher string

const (
	NONE = Cipher("None")
	SALSA20 = Cipher("salsa20")
	AES256CFB = Cipher("aes256cfb")
	AES128CFB = Cipher("aes128cfb")
	SALT = "i'm salt"
)

type EncrytionConfig struct {
	Cipher   Cipher
	Password string
}

func GetKey(k string, kenLen int) []byte {
	pass := pbkdf2.Key([]byte(k), []byte(SALT), 4096, kenLen, sha1.New)
	return pass
}

func NewStreamCipher(config *EncrytionConfig) (CommonCipher, error) {
	switch config.Cipher {
	case SALSA20:
		pass := GetKey(config.Password, 32)
		return cipher.NewSalsa20Stream(pass)
	case AES256CFB:
		pass := GetKey(config.Password, 32)
		return cipher.NewAESStream(pass)
	case AES128CFB:
		pass := GetKey(config.Password, 16)
		return cipher.NewAESStream(pass)
	case NONE:
		return cipher.NewNoneCryptionStream([]byte{})
	default:
		return nil, errors.New("Invalid type!")
	}
}
