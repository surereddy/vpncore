package compress

import (
	"io"
	"github.com/golang/snappy"
)

type SnappyReadWriter struct {
	base io.ReadWriter
	r    io.Reader
	w    io.Writer
}

func NewSnappyReadWriter(base io.ReadWriter) (*SnappyReadWriter, error) {
	return &SnappyReadWriter{
		base: base,
		r: snappy.NewReader(base),
		w: snappy.NewWriter(base),

	}, nil
}

func (this *SnappyReadWriter) Read(data []byte) (int, error) {
	return this.r.Read(data)
}

func (this *SnappyReadWriter) Write(data []byte) (int, error) {
	return this.w.Write(data)
}

func (this *SnappyReadWriter) Close() error {
	this.base = nil

	if c, ok := this.base.(io.Closer); ok {
		return c.(io.Closer).Close()
	}

	return nil
}
