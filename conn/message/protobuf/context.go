package protobuf

import (
	"reflect"
	"github.com/FTwOoO/vpncore/conn"
	"fmt"
)

type ProtobufMessageContext struct {
	ValueToMsgType map[uint16]reflect.Type
	MsgTypeToValue map[reflect.Type]uint16
}

func NewProtobufMessageContext(msgTypes []reflect.Type) (*ProtobufMessageContext, error) {
	self := &ProtobufMessageContext{
		ValueToMsgType:map[uint16]reflect.Type{},
		MsgTypeToValue: map[reflect.Type]uint16{},
	}
	for i, t := range msgTypes {
		if t.Kind() != reflect.Ptr {
			// protobuf's Message type must be pointer
			return nil, fmt.Errorf("Error type that is not pointer %v", t)
		}
		self.ValueToMsgType[uint16(i)] = t
		self.MsgTypeToValue[t] = uint16(i)
	}

	return self, nil
}

func (self *ProtobufMessageContext) PipeMessage(base conn.Message) conn.Message {

	newMsg := &protobufMsg{base:base, ctx:self}
	return newMsg, nil
}

func (self *ProtobufMessageContext) Valid() (bool, error) {
	return len(self.ValueToMsgType) > 0 && len(self.ValueToMsgType) == len(self.MsgTypeToValue)
}


func (this *ProtobufMessageContext) Layer() conn.Layer {
	return conn.APPCATIOIN_LAYER
}


