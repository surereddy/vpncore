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

package protobuf

import (
	"errors"
	"encoding/binary"
)

type protobufPacketHeader struct {
	MessageType uint16
	ContentSize uint16
	Hash        uint64
}

func (d *protobufPacketHeader) HeaderSize() int {
	return binary.Size(d.MessageType) + binary.Size(d.ContentSize) + binary.Size(d.Hash)
}

func (d *protobufPacketHeader) FromBytes(b []byte) (error) {
	if len(b) < d.HeaderSize() {
		return errors.New("Need more data")
	}

	d.MessageType = binary.BigEndian.Uint16(b[:2])
	d.ContentSize = binary.BigEndian.Uint16(b[2:4])
	d.Hash = binary.BigEndian.Uint64(b[4:])
	return nil
}

func (d *protobufPacketHeader) ToBytes() []byte {
	buf := make([]byte, d.HeaderSize())
	binary.BigEndian.PutUint16(buf[:2], d.MessageType)
	binary.BigEndian.PutUint16(buf[2:4], d.ContentSize)
	binary.BigEndian.PutUint64(buf[4:], d.Hash)
	return buf
}

func (d *protobufPacketHeader) ValidateContent(body []byte) error {
	return nil
}
