package msgpack

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *TestMsg) DecodeMsg(dc *msgp.Reader) (err error) {
	var zxvk uint32
	zxvk, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zxvk != 1 {
		err = msgp.ArrayError{Wanted: 1, Got: zxvk}
		return
	}
	z.Data, err = dc.ReadBytes(z.Data)
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *TestMsg) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 1
	err = en.Append(0x91)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.Data)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *TestMsg) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 1
	o = append(o, 0x91)
	o = msgp.AppendBytes(o, z.Data)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TestMsg) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zbzg uint32
	zbzg, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zbzg != 1 {
		err = msgp.ArrayError{Wanted: 1, Got: zbzg}
		return
	}
	z.Data, bts, err = msgp.ReadBytesBytes(bts, z.Data)
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *TestMsg) Msgsize() (s int) {
	s = 1 + msgp.BytesPrefixSize + len(z.Data)
	return
}
