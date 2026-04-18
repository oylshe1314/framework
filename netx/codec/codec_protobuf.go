package codec

import (
	"github.com/oylshe1314/framework/errors"

	"google.golang.org/protobuf/proto"
)

type protobufCodec struct{}

func NewProtobufCodec() Codec {
	return protobufCodec{}
}

func (p protobufCodec) Name() string {
	return "protobuf"
}

func (p protobufCodec) Encode(v any) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, errors.New("not protobuf message")
	}
	return proto.Marshal(msg)
}

func (p protobufCodec) Decode(buf []byte, v any) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return errors.New("not protobuf message")
	}
	return proto.Unmarshal(buf, msg)
}
