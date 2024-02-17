package util

import (
	"fmt"
	"framework/errors"
	json "github.com/json-iterator/go"
	"reflect"
	"strconv"
	"strings"
)

func JsonArrayToBooleans(vs []interface{}) ([]bool, error) {
	var s = make([]bool, len(vs))
	for i, v := range vs {
		rv, ok := v.(bool)
		if !ok {
			return nil, errors.Error("the element of json array was not bool")
		}
		s[i] = rv
	}
	return s, nil
}

func JsonArrayToStrings(vs []interface{}) ([]string, error) {
	var s = make([]string, len(vs))
	for i, v := range vs {
		rv, ok := v.(string)
		if !ok {
			return nil, errors.Error("the element of json array was not string")
		}
		s[i] = rv
	}
	return s, nil
}

func JsonArrayToNumbers[T Number](vs []interface{}, t T) ([]T, error) {
	var s = make([]T, len(vs))
	for i, v := range vs {
		rv, ok := v.(float64)
		if !ok {
			return nil, errors.Error("the element of json array was not number")
		}
		s[i] = T(rv)
	}
	return s, nil
}

func FromJsonObject(obj map[string]interface{}, p interface{}) error {
	if obj == nil || p == nil {
		return nil
	}

	var pv = reflect.ValueOf(p)
	if pv.Kind() != reflect.Pointer {
		return errors.Error("non-pointer")
	}

	pv = pv.Elem()
	if pv.Kind() != reflect.Struct {
		return errors.Error("non-structure")
	}

	return fromJsonObject(obj, pv)
}

func newConvertError(s, t reflect.Type, err ...error) error {
	if len(err) == 0 {
		return errors.Error(fmt.Sprintf("cannot be converted from the type %s to %s", s.String(), t.String()))
	} else {
		return errors.Error(fmt.Sprintf("cannot be converted from the type %s to %s, ", s.String(), t.String()), err[0])
	}
}

func fromJsonValue(val interface{}, tv reflect.Value) error {
	if val == nil {
		return nil
	}

	var tt = tv.Type()
	var vt = reflect.TypeOf(val)
	if vt.ConvertibleTo(tt) {
		tv.Set(reflect.ValueOf(val).Convert(tv.Type()))
		return nil
	} else {
		switch tt.Kind() {
		case reflect.Pointer:
			if tv.IsNil() {
				var pv = reflect.New(tt.Elem())
				var err = fromJsonValue(val, pv.Elem())
				if err != nil {
					return err
				}
				tv.Set(pv)
			} else {
				var err = fromJsonValue(val, tv.Elem())
				if err != nil {
					return err
				}
			}
			return nil
		case reflect.Array, reflect.Slice:
			if vt.Kind() != reflect.Slice {
				return newConvertError(vt, tt)
			}

			ary, ok := val.([]interface{})
			if !ok {
				return newConvertError(vt, tt)
			}

			var err = fromJsonArray(ary, tv)
			if err != nil {
				return err
			}
			return nil
		case reflect.Struct:
			if vt.Kind() != reflect.Map {
				return newConvertError(vt, tt)
			}

			obj, ok := val.(map[string]interface{})
			if !ok {
				return newConvertError(vt, tt)
			}

			var err = fromJsonObject(obj, tv)
			if err != nil {
				return err
			}
			return nil
		case reflect.Map:
			if tt.Key().Kind() != reflect.String {
				return newConvertError(vt, tt)
			}

			if vt.Kind() != reflect.Map {
				return newConvertError(vt, tt)
			}

			obj, ok := val.(map[string]interface{})
			if !ok {
				return newConvertError(vt, tt)
			}
			if tv.IsNil() {
				tv.Set(reflect.MakeMap(tv.Type()))
			}

			if tt.Elem().Kind() != reflect.Interface {
				for jok, jov := range obj {
					tv.SetMapIndex(reflect.ValueOf(jok), reflect.ValueOf(jov))
				}
			} else {
				var tet = tt.Elem()
				for mk, mv := range obj {
					var tev = reflect.Zero(tet)
					var err = fromJsonValue(mv, tev)
					if err != nil {
						return err
					}
					tv.SetMapIndex(reflect.ValueOf(mk), tev)
				}
			}
			return nil
		default:
			switch rlv := val.(type) {
			case bool:
				if tt.Kind() == reflect.String {
					tv.SetString(strconv.FormatBool(rlv))
					return nil
				}
			case float64:
				if tt.Kind() == reflect.String {
					if rlv-float64(int64(rlv)) > 0 {
						tv.SetString(strconv.FormatFloat(rlv, 'f', -1, 64))
					} else {
						tv.SetString(strconv.FormatInt(int64(rlv), 10))
					}
					return nil
				}
			case string:
				switch tt.Kind() {
				case reflect.Bool:
					var b bool
					b, err := strconv.ParseBool(rlv)
					if err != nil {
						return newConvertError(vt, tt, err)
					}
					tv.SetBool(b)
					return nil
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					var n int64
					n, err := strconv.ParseInt(rlv, 10, 64)
					if err != nil {
						return newConvertError(vt, tt, err)
					}
					tv.SetInt(n)
					return nil
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					var u uint64
					u, err := strconv.ParseUint(rlv, 10, 64)
					if err != nil {
						return newConvertError(vt, tt, err)
					}
					tv.SetUint(u)
					return nil
				case reflect.Float32, reflect.Float64:
					var f float64
					f, err := strconv.ParseFloat(rlv, 64)
					if err != nil {
						return newConvertError(vt, tt, err)
					}
					tv.SetFloat(f)
					return nil
				}
			}
			return newConvertError(vt, tt)
		}
	}
}

func fromJsonArray(ary []interface{}, sv reflect.Value) error {
	if ary == nil {
		return nil
	}

	switch sv.Kind() {
	case reflect.Array:
		for i, av := range ary {
			var err = fromJsonValue(av, sv.Index(i))
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		var nsv = reflect.MakeSlice(sv.Type(), len(ary), cap(ary))
		for i, av := range ary {
			var err = fromJsonValue(av, nsv.Index(i))
			if err != nil {
				return err
			}
		}
		sv.Set(nsv)
	}
	return nil
}

func fromJsonObject(obj map[string]interface{}, tv reflect.Value) error {
	if obj == nil {
		return nil
	}

	var tt = tv.Type()
	var fn = tt.NumField()
	for i := 0; i < fn; i++ {
		var sf = tt.Field(i)
		if sf.Anonymous {
			var fv = tv.Field(i)
			if sf.Type.Kind() == reflect.Pointer {
				if fv.IsNil() {
					var fpv = reflect.New(sf.Type.Elem())
					var err = fromJsonObject(obj, fpv.Elem())
					if err != nil {
						return err
					}
					fv.Set(fpv)
				} else {
					var err = fromJsonObject(obj, fv.Elem())
					if err != nil {
						return err
					}
				}
			} else {
				var err = fromJsonObject(obj, fv)
				if err != nil {
					return err
				}
			}
		} else {
			//var omitempty = false
			var name = sf.Tag.Get("json")
			if name == "-" {
				continue
			}

			if name == "" {
				name = sf.Name
			} else {
				var ss = strings.Split(name, ",")
				name = ss[0]
				if len(ss) > 1 {
					if name == "" {
						name = sf.Name
					}
					//if ss[1] == "omitempty" {
					//	omitempty = true
					//}
				}
			}

			var mv = obj[name]
			if mv == nil {
				//ignore empty forever while deserializing
				continue
			}

			var err = fromJsonValue(mv, tv.Field(i))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewReflectValueFromJson(val interface{}, tt reflect.Type) (reflect.Value, error) {
	if val == nil {
		return reflect.Zero(tt), nil
	}

	var isPointer = tt.Kind() == reflect.Pointer
	if isPointer {
		tt = tt.Elem()
	}

	var tv = reflect.New(tt)
	var err = fromJsonValue(val, tv.Elem())
	if err != nil {
		return reflect.Value{}, err
	}
	if !isPointer {
		tv = tv.Elem()
	}

	return tv, nil
}

func ToJsonString(v any) string {
	if v == nil {
		return "<nil>"
	}

	buf, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(buf)
}
