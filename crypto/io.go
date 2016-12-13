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
	"io"
)

type CryptionReadWriter struct {
	stream CommonCipher
	base   io.ReadWriter
}

func NewCryptionReadWriter(base io.ReadWriter, stream CommonCipher) (*CryptionReadWriter, error) {
	if base == nil || stream == nil {
		return nil, ErrArgs
	}

	return &CryptionReadWriter{
		stream: stream,
		base: base,
	}, nil
}

func (this *CryptionReadWriter) Read(data []byte) (int, error) {

	nBytes, err := this.base.Read(data)
	if nBytes > 0 {
		this.stream.Decrypt(data[:nBytes], data[:nBytes])
	}
	return nBytes, err
}

func (this *CryptionReadWriter) Write(data []byte) (int, error) {

	//TODO: use recycle buffer for buf
	//why just copy data to data? because data is client's valid data,
	//dont change it
	buf := make([]byte, len(data))
	this.stream.Encrypt(buf, data)
	return this.base.Write(buf)
}

func (this *CryptionReadWriter) Close() error {
	this.base = nil
	this.stream = nil

	if c, ok := this.base.(io.Closer); ok {
		return c.(io.Closer).Close()
	}

	return nil
}
