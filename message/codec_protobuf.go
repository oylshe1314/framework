package message

type protobufCodec struct {
	base[*protobufCodec]
}

func NewProtobufCodec() Codec {
	return &protobufCodec{}
}

func (p *protobufCodec) encode(msg interface{}) ([]byte, error) {
	panic("implement me")
}

func (p *protobufCodec) decode(buf []byte, msg interface{}) error {
	panic("implement me")
}
