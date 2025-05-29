package message

import "github.com/oylshe1314/framework/errors"

type Codec interface {
	Encode(msg interface{}) ([]byte, error)
	Decode(buf []byte, msg interface{}) error
}

type codec interface {
	encode(msg interface{}) ([]byte, error)
	decode(buf []byte, msg interface{}) error
}

type baseCodec[c codec] struct {
	c c
}

func (this *baseCodec[c]) Encode(msg interface{}) ([]byte, error) {
	switch v := msg.(type) {
	case nil:
		return nil, nil
	case []byte:
		return v, nil
	default:
		return c.encode(this.c, msg)
	}
}

func (this *baseCodec[c]) Decode(buf []byte, msg interface{}) error {
	switch v := msg.(type) {
	case nil:
		return nil
	case []byte:
		copy(v, buf)
	case *[]byte:
		*v = buf
	default:
		return c.decode(this.c, buf, msg)
	}
	return nil
}

type stringCodec struct {
	baseCodec[*stringCodec]
}

func (*stringCodec) encode(msg interface{}) ([]byte, error) {
	switch p := msg.(type) {
	case string:
		return []byte(p), nil
	case *string:
		return []byte(*p), nil
	default:
		return nil, errors.Error("non-string")
	}
}

func (*stringCodec) decode(buf []byte, msg interface{}) error {
	p, ok := msg.(*string)
	if !ok {
		return errors.Error("non-string-pointer")
	}
	*p = string(buf)
	return nil
}

var DefaultCodec Codec = &stringCodec{}
