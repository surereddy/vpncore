package cipher

type NoneCryptionStream struct {
	xortbl []byte
}

func NewNoneCryptionStream(key []byte) (StreamCipher, error) {
	return new(NoneCryptionStream), nil
}

func (c *NoneCryptionStream) Encrypt(dst, src []byte) {
	copy(dst, src)
}

func (c *NoneCryptionStream) Decrypt(dst, src []byte) {
	copy(dst, src)
}
