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
	"io"
	"github.com/golang/protobuf/proto"
	"reflect"
	"fmt"
)

type protobufCodec struct {
	headBuf        []byte
	bodyBuf        []byte
	rw             io.ReadWriter

	maxRecv        int
	maxSend        int
	valueToMsgType map[uint16]reflect.Type
	msgTypeToValue map[reflect.Type]uint16
}

func (c *protobufCodec) Receive() (interface{}, error) {
	if _, err := io.ReadFull(c.rw, c.headBuf); err != nil {
		return nil, err
	}
	header, err := c.DecodeHeader(c.headBuf)
	if err != nil {
		return nil, err

	}
	size := header.ContentSize

	if int(size) > c.maxRecv {
		return nil, errors.New("Too Large Packet")
	}
	if cap(c.bodyBuf) < int(size) {
		c.bodyBuf = make([]byte, size, size + 128)
	}
	buff := c.bodyBuf[:size]
	if _, err := io.ReadFull(c.rw, buff); err != nil {
		return nil, err
	}

	msg, err := c.DecodeBody(buff, header)
	return msg, err
}

func (c *protobufCodec) Send(msg interface{}) error {
	switch msg.(type) {
	case proto.Message:
		all, err := c.EncodeMessage(msg.(proto.Message))
		if err != nil {
			return err
		}

		_, err = c.rw.Write(all)
		if err != nil {
			return err
		}

		return nil
	default:
		return errors.New("Type is not valid")
	}

}

func (c *protobufCodec) Close() error {
	if closer, ok := c.rw.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (d *protobufCodec) EncodeMessage(msg proto.Message) (packet []byte, err error) {
	body, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	h := new(protobufPacketHeader)
	if _, ok := d.msgTypeToValue[reflect.TypeOf(msg)]; !ok {
		fmt.Printf("No message type for this type %s", reflect.TypeOf(msg))
		return nil, fmt.Errorf("No message type for this type %s", reflect.TypeOf(msg))
	}

	h.MessageType = d.msgTypeToValue[reflect.TypeOf(msg)]
	h.ContentSize = uint16(len(body))
	h.Hash = 0

	if err != nil {
		return
	}

	header := h.ToBytes()
	packet = make([]byte, len(header) + len(body))
	copy(packet, header)
	copy(packet[len(header):], body)
	return
}

func (d *protobufCodec) DecodeHeader(header []byte) (h *protobufPacketHeader, err error) {

	h = new(protobufPacketHeader)
	err = h.FromBytes(header)
	return
}

func (d *protobufCodec) DecodeBody(body []byte, h *protobufPacketHeader) (msg proto.Message, err error) {
	if len(body) != int(h.ContentSize) {
		return nil, errors.New("Content size dont match")
	}

	if err := h.ValidateContent(body); err != nil {
		return nil, err
	}

	if _, ok := d.valueToMsgType[h.MessageType]; !ok {
		return nil, errors.New("Invalid type")
	}

	T := d.valueToMsgType[h.MessageType]
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

