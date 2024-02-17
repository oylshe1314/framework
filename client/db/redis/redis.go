package redis

import (
	"context"
	"fmt"
	"framework/errors"
	"reflect"
	"strconv"
)

type Strings []string
type StringMap map[string]string

func (sm StringMap) Encode(v interface{}) error {
	if sm == nil || v == nil {
		return errors.Error("encode nil-pointer")
	}

	var vt = reflect.TypeOf(v)
	var vv = reflect.ValueOf(v)
	if vt.Kind() == reflect.Pointer {
		vt = vt.Elem()
		vv = vv.Elem()
	}

	if vt.Kind() != reflect.Struct {
		return errors.Error("decode non-structure")
	}

	var fn = vt.NumField()
	for i := 0; i < fn; i++ {
		var sf = vt.Field(i)
		if !sf.IsExported() {
			continue
		}

		var ft = sf.Type
		var fv = vv.Field(i)
		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
			fv = fv.Elem()
		}

		var name = sf.Tag.Get("redis")
		if name == "-" {
			continue
		}

		if len(name) == 0 {
			name = sf.Name
		}

		switch ft.Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			sm[name] = fmt.Sprint(fv.Interface())
		case reflect.String:
			sm[name] = fv.String()
		default:
			return errors.Errorf("field '%s' type '%s' cannot parse encode", name, ft.String())
		}
	}

	return nil
}

func (sm StringMap) Decode(v interface{}) error {
	var vt = reflect.TypeOf(v)
	if vt.Kind() != reflect.Pointer {
		return errors.Error("decode non-pointer")
	}

	vt = vt.Elem()
	if vt.Kind() != reflect.Struct {
		return errors.Error("decode non-structure")
	}

	var vv = reflect.ValueOf(v).Elem()
	var fn = vt.NumField()
	for i := 0; i < fn; i++ {
		var sf = vt.Field(i)
		if !sf.IsExported() {
			continue
		}

		var ft = sf.Type
		var fv = vv.Field(i)
		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
			fv = fv.Elem()
		}

		var name = sf.Tag.Get("redis")
		if name == "-" {
			continue
		}

		if len(name) == 0 {
			name = sf.Name
		}

		val, ok := sm[name]
		if !ok {
			continue
		}

		switch ft.Kind() {
		case reflect.Bool:
			b, err := strconv.ParseBool(val)
			if err != nil {
				return errors.Errorf("field '%s' cannot parse bool, %v", name, err)
			}
			fv.SetBool(b)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return errors.Errorf("field '%s' cannot parse int, %v", name, err)
			}
			fv.SetInt(n)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			n, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return errors.Errorf("field '%s' cannot parse uint, %v", name, err)
			}
			fv.SetUint(n)
		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return errors.Errorf("field '%s' cannot parse float, %v", name, err)
			}
			fv.SetFloat(f)
		case reflect.String:
			fv.SetString(val)
		}
	}
	return nil
}

type Redis interface {
	Close() error
	Exec(ctx context.Context, cmd string, args ...interface{}) error
	String(ctx context.Context, cmd string, args ...interface{}) (string, error)
	Strings(ctx context.Context, cmd string, args ...interface{}) (Strings, error)
	StringMap(ctx context.Context, cmd string, args ...interface{}) (StringMap, error)
	Subscribe(ctx context.Context) *SubConn
}
