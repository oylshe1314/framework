package server

import (
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/option"
	"github.com/oylshe1314/framework/util"
	"github.com/oylshe1314/framework/util/jsonx"
	"reflect"
	"strings"
)

func SetOption(svr Server, opt option.Option) ([]Component[Server], error) {
	return serverOption(opt).setOption(svr)
}

type serverOption option.Option

func (opt serverOption) setOption(svr Server) ([]Component[Server], error) {
	var st = reflect.TypeOf(svr)
	if st.Kind() != reflect.Pointer || st.Elem().Kind() != reflect.Struct {
		return nil, errors.New("server should be a struct pointer")
	}

	return opt.reflectSetOption(st.Elem().Name(), reflect.ValueOf(svr), st)
}

func (opt serverOption) reflectSetOption(name string, sv reflect.Value, st reflect.Type) ([]Component[Server], error) {
	var err error
	var components []Component[Server]

	var sit = reflect.TypeOf((*Server)(nil)).Elem()

	st = st.Elem()
	sv = sv.Elem()

	for i := range st.NumField() {
		var sf = st.Field(i)
		if !sf.IsExported() {
			continue
		}

		var ft = sf.Type
		if ft.Kind() != reflect.Pointer {
			ft = reflect.PointerTo(ft)
		}

		if ft.Implements(sit) {
			var tags = strings.Split(sf.Tag.Get("server"), ",")

			var svrName = tags[0]
			if svrName == "" {
				svrName = util.LowerCamelCase(sf.Name)
			}

			var val = option.Option(opt).Get(svrName)
			if val == nil {
				continue
			}

			tmp, ok := val.(map[string]any)
			if !ok {
				return nil, errors.Errorf("can not set option '%s' of %s to server '%s' ", name, reflect.TypeOf(tmp).String(), sf.Name)
			}

			var fv = sv.Field(i)
			if sf.Type.Kind() == reflect.Pointer {
				if fv.IsNil() {
					return nil, errors.Errorf("the pointer of the server '%s' is nil", sf.Name)
				}
			} else {
				fv = fv.Addr()
			}

			var subComponents []Component[Server]
			subComponents, err = (serverOption(tmp)).reflectSetOption(svrName, fv, ft)
			if err != nil {
				return nil, err
			}

			components = append(components, subComponents...)
		}
	}

	for i := range st.NumMethod() {
		var sm = st.Method(i)
		if !sm.IsExported() {
			continue
		}

		if sm.Name == "SetOption" {
			var mt = sm.Type
			if mt.NumIn() != 1 {
				return nil, errors.Errorf("invalid parameter of the server '%s' option set funcation", st.Name())
			}

			var at = mt.In(0)
			if at.Kind() != reflect.Struct || at.Elem().Kind() != reflect.Struct {
				return nil, errors.Errorf("the parameter of the server '%s' option set funcation should be a struct or struct pointer", st.Name())
			}

			if at.Kind() == reflect.Pointer {
				var av = reflect.New(at.Elem())
				err = jsonx.Object(opt).ToStruct(av.Interface())
				if err != nil {
					return nil, err
				}

				sv.Method(i).Call([]reflect.Value{av})
			} else {
				var av = reflect.New(at)
				err = jsonx.Object(opt).ToStruct(av.Interface())
				if err != nil {
					return nil, err
				}
				sv.Method(i).Call([]reflect.Value{av.Elem()})
			}
		}
	}

	for st.Kind() == reflect.Pointer {
		st = st.Elem()
	}

	components = append(components, NewServerComponent(name, sv.Interface().(Server)))
	return components, nil
}
