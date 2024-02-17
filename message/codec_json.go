package message

import (
	json "github.com/json-iterator/go"
)

type jsonCodec struct {
	base[*jsonCodec]
}

func NewJsonCodec() Codec {
	return &jsonCodec{}
}

func (*jsonCodec) encode(msg interface{}) ([]byte, error) {
	return json.Marshal(msg)
}

func (*jsonCodec) decode(buf []byte, msg interface{}) error {
	return json.Unmarshal(buf, msg)
}
