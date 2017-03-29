package compress

import (
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
	return NewSnappyReadWriter(base)
}
