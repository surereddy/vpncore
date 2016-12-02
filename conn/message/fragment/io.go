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

	closedM bool
	closedC chan struct{}
}

func NewFragmentIO(base conn.StreamIO) (*FragmentIO, error) {
	f := &FragmentIO{
		base: base,
		closedM:false,
		closedC:make(chan struct{}),
	}

	return f, nil
}

func (this *FragmentIO) Read() (buf []byte, err error) {
	var lengthBytes [2]byte
	if _, err = io.ReadFull(this.base, lengthBytes[:]); err != nil {
		return
	}

	length := binary.BigEndian.Uint16(lengthBytes[:])
	buf = make([]byte, length)

	if _, err = io.ReadFull(this.base, buf[:length]); err != nil {
		return
	}

	return
}

func (this *FragmentIO) Write(b []byte) error {
	if b == nil || len(b) == 0 {
		return nil
	}

	var lengthBytes [2]byte
	binary.BigEndian.PutUint16(lengthBytes[:], uint16(len(b)))

	_, err := this.base.Write(lengthBytes[:])
	if err != nil {
		return err
	}

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