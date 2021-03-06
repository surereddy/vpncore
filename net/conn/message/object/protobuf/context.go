package protobuf

import (
	"reflect"
	"github.com/FTwOoO/vpncore/net/conn"
	"fmt"
	"github.com/golang/protobuf/proto"
)

var _ conn.MessageToObjectContext = new(ProtobufMessageContext)

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

func (self *ProtobufMessageContext) Valid() (bool, error) {
	res := len(self.ValueToMsgType) > 0 && len(self.ValueToMsgType) == len(self.MsgTypeToValue)
	return res, nil
}

func (this *ProtobufMessageContext) Layer() conn.Layer {
	return conn.APPCATIOIN_LAYER
}

func (this *ProtobufMessageContext) Encode(obj interface{}) ([]byte, error) {

	if _, ok := obj.(proto.Message); !ok {
		return nil, conn.ErrUnsupportType
	}

	msg := protobufMsg{Content:obj.(proto.Message), ctx:this}
	return msg.Marshal()
}

func (this *ProtobufMessageContext) Decode(packet []byte) (interface{}, error) {
	msg := protobufMsg{Content:nil, ctx:this}
	_, err := msg.Unmarshal(packet)

	if err != nil {
		return nil, err
	}

	return msg.Content, nil
}

