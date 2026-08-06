package jsonx

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/util"
)

type Object map[string]any

func newReflectConvertError(s, t reflect.Type, err ...error) error {
	if len(err) == 0 {
		return errors.Error(fmt.Sprintf("cannot be converted from the type %s to %s", s.String(), t.String()))
	} else {
		return errors.Error(fmt.Sprintf("cannot be converted from the type %s to %s, ", s.String(), t.String()), err[0])
	}
}

func (obj Object) ToStruct(p any) error {
	var pt = reflect.TypeOf(p)
	if pt.Kind() != reflect.Pointer || pt.Elem().Kind() != reflect.Struct {
		return errors.Error("not a struct pointer")
	}

	var ppv = reflect.New(pt)
	ppv.Elem().Set(reflect.ValueOf(p))

	return obj.reflectToStruct(ppv)
}

func (obj Object) reflectToStruct(ppv reflect.Value) error {
	var pv = ppv.Elem()
	var vt = pv.Type().Elem()

	var vv reflect.Value
	if !pv.IsNil() {
		vv = pv.Elem()
	}

	for i := range vt.NumField() {
		var sf = vt.Field(i)
		if !sf.IsExported() {
			continue
		}

		var tag = sf.Tag.Get("json")
		if tag == "-" {
			continue
		}

		var name string
		var tags []string
		if tag == "" {
			name = sf.Name
		} else {
			tags = strings.Split(tag, ",")
			name = tags[0]
			if name == "" {
				name = sf.Name
			}
		}

		if sf.Anonymous {
			var val map[string]any
			if slices.Contains(tags, "inline") {
				val = obj
			} else {
				tmp, ok := obj[name]
				if !ok || tmp == nil {
					continue
				}

				val, ok = tmp.(map[string]any)
				if !ok {
					return newReflectConvertError(reflect.TypeOf(tmp), sf.Type)
				}
			}

			if sf.Type.Kind() == reflect.Pointer {
				if pv.IsNil() {
					var fpp = reflect.New(sf.Type)
					var err = Object(val).reflectToStruct(fpp)
					if err != nil {
						return err
					}
					if !fpp.Elem().IsNil() {
						pv.Set(reflect.New(vt))
						vv = pv.Elem()
						vv.Field(i).Set(fpp.Elem())
					}
				} else {
					var err = Object(val).reflectToStruct(vv.Field(i).Addr())
					if err != nil {
						return err
					}
				}
			} else {
				var fpp = reflect.New(reflect.PointerTo(sf.Type))
				if pv.IsNil() {
					var err = Object(val).reflectToStruct(fpp)
					if err != nil {
						return err
					}
					if !fpp.Elem().IsNil() {
						pv.Set(reflect.New(vt))
						vv = pv.Elem()
						vv.Field(i).Set(fpp.Elem().Elem())
					}
				} else {
					fpp.Elem().Set(vv.Field(i).Addr())
					var err = Object(val).reflectToStruct(fpp)
					if err != nil {
						return err
					}
				}
			}
		} else {
			val, ok := obj[name]
			if !ok || val == nil {
				continue
			}

			if pv.IsNil() {
				pv.Set(reflect.New(vt))
				vv = pv.Elem()
			}

			var err = reflectAnyToValue(val, vv.Field(i), sf.Type)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func AnyToValue(val any, ptr any) error {
	var pv = reflect.ValueOf(ptr).Elem()
	return reflectAnyToValue(val, pv, pv.Type())
}

func reflectAnyToValue(val any, vv reflect.Value, ft reflect.Type) error {
	var vt = vv.Type()
	if reflect.TypeOf(val).ConvertibleTo(vt) {
		vv.Set(reflect.ValueOf(val).Convert(vt))
	} else {
		if ft == nil {
			ft = vt
		}

		switch vt.Kind() {
		case reflect.Bool:
			switch tv := val.(type) {
			case bool:
				vv.SetBool(tv)
			case float64:
				vv.SetBool(tv != 0)
			case string:
				b, err := strconv.ParseBool(tv)
				if err != nil {
					return newReflectConvertError(reflect.TypeOf(tv), ft, err)
				}
				vv.SetBool(b)
			default:
				return newReflectConvertError(reflect.TypeOf(tv), ft)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			switch tv := val.(type) {
			case bool:
				vv.SetInt(util.If[int64](tv, 1, 0))
			case float64:
				vv.SetInt(int64(tv))
			case string:
				n, err := strconv.ParseInt(tv, 10, 64)
				if err != nil {
					return newReflectConvertError(reflect.TypeOf(tv), ft, err)
				}
				vv.SetInt(n)
			default:
				return newReflectConvertError(reflect.TypeOf(tv), ft)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			switch tv := val.(type) {
			case bool:
				vv.SetUint(util.If[uint64](tv, 1, 0))
			case float64:
				vv.SetUint(uint64(tv))
			case string:
				u, err := strconv.ParseUint(tv, 10, 64)
				if err != nil {
					return newReflectConvertError(reflect.TypeOf(tv), ft, err)
				}
				vv.SetUint(u)
			default:
				return newReflectConvertError(reflect.TypeOf(tv), ft)
			}
		case reflect.Float32, reflect.Float64:
			switch tv := val.(type) {
			case bool:
				vv.SetFloat(util.If[float64](tv, 1.0, 0.0))
			case float64:
				vv.SetFloat(tv)
			case string:
				f, err := strconv.ParseFloat(tv, 64)
				if err != nil {
					return newReflectConvertError(reflect.TypeOf(tv), ft, err)
				}
				vv.SetFloat(f)
			default:
				return newReflectConvertError(reflect.TypeOf(tv), ft)
			}
		case reflect.Array:
			switch tv := val.(type) {
			case []any:
				var length = vv.Len()
				if length > len(tv) {
					length = len(tv)
				}
				for i := range length {
					var iv = vv.Index(i)
					var err = reflectAnyToValue(tv[i], iv, nil)
					if err != nil {
						return err
					}
				}
			default:
				return newReflectConvertError(reflect.TypeOf(val), ft)
			}
		case reflect.Interface:
			vv.Set(reflect.ValueOf(val))
		case reflect.Map:
			switch tv := val.(type) {
			case map[string]any:
				if vv.IsNil() {
					vv.Set(reflect.MakeMap(vt))
				}
				for mk, mv := range tv {
					var err error
					var rk = reflect.New(vt.Key()).Elem()
					err = reflectAnyToValue(mk, rk, nil)
					if err != nil {
						return err
					}

					var rv = reflect.New(vt.Elem()).Elem()
					err = reflectAnyToValue(mv, rv, nil)
					if err != nil {
						return err
					}

					vv.SetMapIndex(rk, rv)
				}
			default:
				return newReflectConvertError(reflect.TypeOf(val), ft)
			}
		case reflect.Pointer:
			if vv.IsNil() {
				if vt.Elem().Kind() == reflect.Struct {
					tmp, ok := val.(map[string]any)
					if !ok {
						return newReflectConvertError(reflect.TypeOf(val), ft)
					}
					var ppv = reflect.New(vt)
					var err = Object(tmp).reflectToStruct(ppv)
					if err != nil {
						return err
					}
					if !ppv.Elem().IsNil() {
						vv.Set(ppv.Elem())
					}
				} else {
					var pv = reflect.New(vt.Elem())
					var err = reflectAnyToValue(val, pv.Elem(), ft)
					if err != nil {
						return err
					}
					vv.Set(pv)
				}
			} else {
				var err = reflectAnyToValue(val, vv.Elem(), ft)
				if err != nil {
					return err
				}
			}
		case reflect.Slice:
			switch tv := val.(type) {
			case []any:
				if vv.IsNil() {
					var length = len(tv)
					var sv = reflect.MakeSlice(vt, length, cap(tv))
					for i := range length {
						var iv = sv.Index(i)
						var err = reflectAnyToValue(tv[i], iv, nil)
						if err != nil {
							return err
						}
					}
					vv.Set(sv)
				} else {
					var length = vv.Len()
					if length > len(tv) {
						length = len(tv)
					}
					for i := range length {
						var iv = vv.Index(i)
						var err = reflectAnyToValue(tv[i], iv, nil)
						if err != nil {
							return err
						}
					}
				}
			default:
				return newReflectConvertError(reflect.TypeOf(val), ft)
			}
		case reflect.String:
			switch sv := val.(type) {
			case bool:
				vv.SetString(strconv.FormatBool(sv))
			case float64:
				vv.SetString(strconv.FormatFloat(sv, 'f', -1, 64))
			case string:
				vv.SetString(sv)
			default:
				return newReflectConvertError(reflect.TypeOf(val), ft)
			}
		case reflect.Struct:
			switch tv := val.(type) {
			case map[string]any:
				var ppv = reflect.New(reflect.PointerTo(vt))
				ppv.Elem().Set(vv.Addr())
				var err = Object(tv).reflectToStruct(ppv)
				if err != nil {
					return err
				}
			default:
				return newReflectConvertError(reflect.TypeOf(val), ft)
			}
		default:
			return newReflectConvertError(reflect.TypeOf(val), ft)
		}
	}
	return nil
}
