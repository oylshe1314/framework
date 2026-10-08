package codec

import "encoding/json"

type jsonCodec struct{}

func NewJsonCodec() Codec {
	return jsonCodec{}
}

func (jsonCodec) Name() string {
	return "json"
}

func (jsonCodec) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (jsonCodec) Decode(buf []byte, v any) error {
	return json.Unmarshal(buf, v)
}
