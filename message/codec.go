package message

import "framework/errors"

type Codec interface {
	Encode(msg interface{}) ([]byte, error)
	Decode(buf []byte, msg interface{}) error
}

type codec interface {
	encode(msg interface{}) ([]byte, error)
	decode(buf []byte, msg interface{}) error
}

type base[c codec] struct {
	c c
}

func (this *base[c]) Encode(msg interface{}) ([]byte, error) {
	switch v := msg.(type) {
	case nil:
		return nil, nil
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		return c.encode(this.c, msg)
	}
}

func (this *base[c]) Decode(buf []byte, msg interface{}) error {
	switch v := msg.(type) {
	case nil:
		return nil
	case []byte:
		copy(v, buf)
		return nil
	case *[]byte:
		*v = buf
		return nil
	case *string:
		*v = string(buf)
		return nil
	default:
		return c.decode(this.c, buf, msg)
	}
}

type stringCodec struct {
	base[*stringCodec]
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
