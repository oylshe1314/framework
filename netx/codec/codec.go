package codec

import (
	"github.com/oylshe1314/framework/errors"
)

type Codec interface {
	Name() string

	Encode(any) ([]byte, error)
	Decode([]byte, any) error
}

func NewCodec(name string) (Codec, error) {
	switch name {
	case "string":
		return NewStringCodec(), nil
	case "json":
		return NewJsonCodec(), nil
	case "protobuf":
		return NewProtobufCodec(), nil
	default:
		return nil, errors.Errorf("undefined codec: %s", name)
	}
}

type stringCodec struct{}

func NewStringCodec() Codec {
	return stringCodec{}
}

func (codec stringCodec) Name() string {
	return "string"
}

func (codec stringCodec) Encode(v any) ([]byte, error) {
	switch vt := v.(type) {
	case []byte:
		return vt, nil
	case string:
		return []byte(vt), nil
	case *string:
		return []byte(*vt), nil
	default:
		return nil, errors.Errorf("'%T' cannot encode to []byte", v)
	}
}

func (codec stringCodec) Decode(buf []byte, v any) error {
	switch vt := v.(type) {
	case []byte:
		copy(vt, buf)
	case string:
		return errors.New("not string pointer")
	case *string:
		*vt = string(buf)
	default:
		return errors.Errorf("'%T' cannot decode from []byte", v)
	}
	return nil
}

var (
	DefaultCodec = NewStringCodec()
)
