package message

import (
	"github.com/oylshe1314/framework/errors"
	"google.golang.org/protobuf/proto"
)

type protobufCodec struct {
	baseCodec[*protobufCodec]
}

func NewProtobufCodec() Codec {
	return &protobufCodec{}
}

func (*protobufCodec) encode(msg interface{}) ([]byte, error) {
	p, ok := msg.(proto.Message)
	if !ok {
		return nil, errors.Error("not protobuf message")
	}
	return proto.Marshal(p)
}

func (*protobufCodec) decode(buf []byte, msg interface{}) error {
	switch p := msg.(type) {
	case proto.Message:
		return proto.Unmarshal(buf, p)
	default:
		return errors.Error("not protobuf message")
	}
}
