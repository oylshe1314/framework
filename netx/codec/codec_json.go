package codec

import "encoding/json"

type jsonCodec struct{}

func NewJsonCodec() Codec {
	return jsonCodec{}
}

func (j jsonCodec) Name() string {
	return "json"
}

func (j jsonCodec) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j jsonCodec) Decode(buf []byte, v any) error {
	return json.Unmarshal(buf, v)
}
