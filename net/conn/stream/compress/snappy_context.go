package compress

import (
	"github.com/FTwOoO/vpncore/crypto"
	"github.com/FTwOoO/vpncore/net/conn"
)

var _ conn.StreamContext = new(SnappyCompressionStreamContext)

type SnappyCompressionStreamContext struct{}

func (this *SnappyCompressionStreamContext) Layer() conn.Layer {
	return conn.COMPRESS_LAYER
}

func (this *SnappyCompressionStreamContext) Valid() (bool, error) {
	return true, nil
}

func (this *SnappyCompressionStreamContext) Pipe(base conn.StreamIO) (c conn.StreamIO) {
	cipher, err := crypto.NewStreamCipher(this.EncrytionConfig)
	if err != nil {
		return nil
	}

	c, err = crypto.NewCryptionReadWriter(base, cipher)
	if err != nil {
		return nil
	}

	return

}
