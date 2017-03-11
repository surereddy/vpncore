package crypto

import (
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha1"
	"errors"
	"github.com/FTwOoO/vpncore/crypto/cipher"
)

var (
	ErrArgs = errors.New("Arguments are not valid.")
)

//TODO: More elegant style to import these types from sub modules to
// local namespace?
type CommonCipher cipher.CommonCipher
type StreamCipher cipher.StreamCipher
type BlockCipher cipher.BlockCipher

type StreamCipherName string

const (
	NONE = StreamCipherName("None")
	AES256CFB = StreamCipherName("aes256cfb")
	AES128CFB = StreamCipherName("aes128cfb")

	SALT = "i'm salt"
)

type EncrytionConfig struct {
	Cipher   StreamCipherName
	Password string
}

func GetKey(k string, kenLen int) []byte {
	pass := pbkdf2.Key([]byte(k), []byte(SALT), 4096, kenLen, sha1.New)
	return pass
}

func NewStreamCipher(config *EncrytionConfig) (CommonCipher, error) {
	switch config.Cipher {
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
