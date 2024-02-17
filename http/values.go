package http

import (
	"framework/errors"
	"reflect"
	"strconv"
	"strings"
)

type UrlValue string

func (value UrlValue) read(v reflect.Value) error {
	var rawValue = string(value)
	switch v.Kind() {
	case reflect.String:
		v.SetString(rawValue)
	case reflect.Bool:
		ev, err := strconv.ParseBool(rawValue)
		if err != nil {
			return err
		}
		v.SetBool(ev)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		ev, err := strconv.ParseInt(rawValue, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(ev)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		ev, err := strconv.ParseUint(rawValue, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(ev)
	case reflect.Float32, reflect.Float64:
		ev, err := strconv.ParseFloat(rawValue, 64)
		if err != nil {
			return err
		}
		v.SetFloat(ev)
	default:
		return errors.Errorf("unsupported type '%s'", v.Type().String())
	}
	return nil
}

type UrlValues map[string][]string

func (values UrlValues) read(v interface{}) error {
	var vt = reflect.TypeOf(v)
	if vt.Kind() != reflect.Pointer {
		return errors.Error("read get query: non-pointer")
	}

	vt = vt.Elem()
	if vt.Kind() != reflect.Struct {
		return errors.Error("read get query: non-struct")
	}

	var vv = reflect.ValueOf(v).Elem()
	var fn = vt.NumField()
	for i := 0; i < fn; i++ {
		var sf = vt.Field(i)
		var ft = sf.Type
		var name = sf.Tag.Get("json")
		if name == "-" {
			continue
		}

		if name == "" {
			name = sf.Name
		} else {
			var chi = strings.IndexByte(name, ',')
			if chi > 0 {
				name = name[:chi]
			}
		}

		var value = values[name]
		if len(value) == 0 {
			continue
		}

		var vl = len(value)
		var fv = vv.Field(i)
		switch ft.Kind() {
		case reflect.Slice:
			if fv.IsNil() {
				fv.Set(reflect.MakeSlice(ft, vl, vl))
			}
			fallthrough
		case reflect.Array:
			if vl > fv.Len() {
				vl = fv.Len()
			}
			for fi := 0; fi < vl; fi++ {
				var ev = fv.Index(fi)
				var err = UrlValue(value[fi]).read(ev)
				if err != nil {
					return errors.Errorf("can not read the value '%s' for index %d of the array field '%s', %v", value[fi], fi, sf.Name, err)
				}
			}
		default:
			var err = UrlValue(value[0]).read(fv)
			if err != nil {
				return errors.Errorf("can not read the value '%s' for the field '%s', %v", value[0], sf.Name, err)
			}
		}
	}

	return nil
}
