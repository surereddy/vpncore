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
	"github.com/golang/protobuf/proto"
	"reflect"
	"fmt"
	"errors"
	"encoding/binary"
)

type protobufMsgHeader struct {
	MessageType uint16
	ContentSize uint16
	Hash        uint64
}

func (d *protobufMsgHeader) HeaderSize() int {
	return binary.Size(d.MessageType) + binary.Size(d.ContentSize) + binary.Size(d.Hash)
}

func (d *protobufMsgHeader) Unmarshal(b []byte) (error) {
	if len(b) < d.HeaderSize() {
		return errors.New("Need more data")
	}

	d.MessageType = binary.BigEndian.Uint16(b[:2])
	d.ContentSize = binary.BigEndian.Uint16(b[2:4])
	d.Hash = binary.BigEndian.Uint64(b[4:])
	return nil
}

func (d *protobufMsgHeader) Marshal() ([]byte, error) {
	buf := make([]byte, d.HeaderSize())
	binary.BigEndian.PutUint16(buf[:2], d.MessageType)
	binary.BigEndian.PutUint16(buf[2:4], d.ContentSize)
	binary.BigEndian.PutUint64(buf[4:], d.Hash)
	return buf, nil
}

func (d *protobufMsgHeader) ValidateContent(body []byte) error {
	return nil
}

type protobufMsg struct {
	Content proto.Message
	ctx     *ProtobufMessageContext
}

func (d *protobufMsg) Unmarshal(buf []byte) (n int, err error) {
	header, err := d.decodeHeader(buf)
	if err != nil {
		return
	}

	buf = buf[header.HeaderSize():]
	if len(buf) < header.HeaderSize() {
		return 0, errors.New("Body content not enough!")
	}

	msg, err := d.decodeBody(buf, header)
	if err != nil {
		return
	} else {

		d.Content = msg
		n = header.HeaderSize() + header.HeaderSize()
		return
	}
}

func (d *protobufMsg) Marshal() (packet []byte, err error) {
	body, err := proto.Marshal(d.Content)
	if err != nil {
		return nil, err
	}

	h := new(protobufMsgHeader)
	if _, ok := d.ctx.MsgTypeToValue[reflect.TypeOf(d.Content)]; !ok {
		fmt.Printf("No message type for this type %s", reflect.TypeOf(d.Content))
		return nil, fmt.Errorf("No message type for this type %s", reflect.TypeOf(d.Content))
	}

	h.MessageType = d.ctx.MsgTypeToValue[reflect.TypeOf(d.Content)]
	h.ContentSize = uint16(len(body))
	h.Hash = 0

	if err != nil {
		return
	}

	header, _ := h.Marshal()
	packet = make([]byte, len(header) + len(body))
	copy(packet, header)
	copy(packet[len(header):], body)

	return
}

func (d *protobufMsg) decodeHeader(header []byte) (h *protobufMsgHeader, err error) {
	h = new(protobufMsgHeader)
	err = h.Unmarshal(header)
	return
}

func (d *protobufMsg) decodeBody(body []byte, h *protobufMsgHeader) (msg proto.Message, err error) {
	if len(body) != int(h.ContentSize) {
		return nil, errors.New("Content size dont match")
	}

	if err := h.ValidateContent(body); err != nil {
		return nil, err
	}

	if _, ok := d.ctx.ValueToMsgType[h.MessageType]; !ok {
		return nil, errors.New("Invalid type")
	}

	T := d.ctx.ValueToMsgType[h.MessageType]
	if T == nil {
		return nil, errors.New("No Message configured for this type")
	}

	msg = reflect.New(T.Elem()).Interface().(proto.Message)
	err = proto.Unmarshal(body, msg)
	if err != nil {
		return nil, err
	}
	return
}
