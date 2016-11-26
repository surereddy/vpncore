package protobuf

import (
	"io"
	"math"
	"reflect"
	"github.com/FTwOoO/vpncore/conn"
)

type ProtobufMessageContext struct {
	msgTypes []reflect.Type
	rw       io.ReadWriter
}

func NewProtobufMessageContext(msgTypes []reflect.Type, rw io.ReadWriter) *ProtobufMessageContext {
	return &ProtobufMessageContext{msgTypes:msgTypes, rw:rw}
}

func (self *ProtobufMessageContext)NewCodec(_ conn.MessageConn) (conn.MessageConn, error) {

	codec := &protobufCodec{
		rw: self.rw,
		headBuf: make([]byte, new(protobufPacketHeader).HeaderSize()),
		maxRecv : math.MaxUint16,
		maxSend : math.MaxUint16,
		valueToMsgType : map[uint16]reflect.Type{},
		msgTypeToValue : map[reflect.Type]uint16{},
	}

	for i, t := range self.msgTypes {
		if t.Kind() != reflect.Ptr {
			// protobuf's Message type must be pointer
			return nil
		}
		codec.valueToMsgType[uint16(i)] = t

		codec.msgTypeToValue[t] = uint16(i)
	}

	return codec, nil
}

func (self *ProtobufMessageContext) Valid() (bool, error) {
	return len(self.msgTypes) > 0 && self.rw != nil
}

