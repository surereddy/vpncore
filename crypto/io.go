/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * Author: FTwOoO <booobooob@gmail.com>
 */

package crypto

import (
	"github.com/FTwOoO/vpncore/common"
	"io"
)

type CryptionReadWriter struct {
	stream CommonCipher
	base   io.ReadWriter
}

func NewCryptionReadWriter(base io.ReadWriter, stream CommonCipher) *CryptionReadWriter {
	return &CryptionReadWriter{
		stream: stream,
		base: base,
	}
}

func (this *CryptionReadWriter) Read(data []byte) (int, error) {
	if this.base == nil {
		return 0, common.ErrObjectNotFound
	}
	nBytes, err := this.base.Read(data)
	if nBytes > 0 {
		this.stream.Decrypt(data[:nBytes], data[:nBytes])
	}
	return nBytes, err
}

func (this *CryptionReadWriter) Write(data []byte) (int, error) {
	if this.base == nil {
		return 0, common.ErrObjectNotFound
	}

	//TODO: use recycle buffer for buf
	//why just copy data to data? because data is client's valid data,
	//dont change it
	buf := make([]byte, len(data))
	this.stream.Encrypt(buf, data)
	return this.base.Write(buf)
}

func (this *CryptionReadWriter) Close() error {
	if c, ok := this.base.(io.Closer); ok {
		c.(io.Closer).Close()
	}

	this.base = nil
	this.stream = nil
}
