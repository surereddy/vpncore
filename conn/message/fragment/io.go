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

package fragment

import (
	"github.com/FTwOoO/vpncore/conn"
	"io"
	"encoding/binary"
)

type FragmentIO struct {
	base    conn.StreamIO
	buf     []byte

	closedM bool
	closedC chan struct{}
}

func NewFragmentIO(base conn.StreamIO) (*FragmentIO, error) {
	f := &FragmentIO{
		base: base,
		closedM:false,
		closedC:make(chan struct{}),
		buf: make([]byte, 0x10000),
	}

	return f, nil
}

func (this *FragmentIO) Read() (buf []byte, err error) {
	var length uint16
	var lengthBytes [2]byte
	if _, err = io.ReadFull(this.base, lengthBytes[:]); err != nil {
		return
	}

	binary.BigEndian.PutUint16(lengthBytes[:], length)
	if _, err = io.ReadFull(this.base, this.buf[:length]); err != nil {
		return
	}

	buf = make([]byte, length)
	copy(buf, this.buf[:length])
	return
}

func (this *FragmentIO) Write(b []byte) error {

	for {
		if len(b) <= 0 {
			break
		}

		n, err := this.base.Write(b)
		if err != nil {
			return err
		}
		b = b[n:]
	}

	return nil
}

func (this *FragmentIO) Close() error {
	return this.base.Close()
}