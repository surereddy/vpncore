package protobuf

import (
	"math"
	"reflect"
	"github.com/FTwOoO/vpncore/conn"
)

type ProtobufMessageContext struct {
	msgTypes []reflect.Type
}

func NewProtobufMessageContext(msgTypes []reflect.Type) *ProtobufMessageContext {
	return &ProtobufMessageContext{msgTypes:msgTypes}
}

func (self *ProtobufMessageContext) NewPipe(base conn.MessagegReadWriteCloser) conn.MessagegReadWriteCloser {

	codec := &protobufCodec{
		rw: base,
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
	return len(self.msgTypes) > 0
}

