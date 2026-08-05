package framework

import (
	"reflect"
	"strings"

	"github.com/oylshe1314/framework/client"
	"github.com/oylshe1314/framework/component"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/option"
	"github.com/oylshe1314/framework/server"
	"github.com/oylshe1314/framework/util"
	"github.com/oylshe1314/framework/util/jsonx"
)

type Option option.Option

func (opt Option) setOption(optional any) ([]component.Component, error) {
	var st = reflect.TypeOf(optional)
	if st.Kind() != reflect.Pointer || st.Elem().Kind() != reflect.Struct {
		return nil, errors.New("server should be a struct pointer")
	}

	return opt.reflectSetOption(util.LowerCamelCase(st.Elem().Name()), reflect.ValueOf(optional), st)
}

func (opt Option) reflectSetOption(name string, spv reflect.Value, spt reflect.Type) ([]component.Component, error) {
	var err error
	var cs []component.Component

	var serverType = reflect.TypeOf((*server.Server)(nil)).Elem()
	var clientType = reflect.TypeOf((*client.Client)(nil)).Elem()

	var cc component.Component
	switch {
	case spt.Implements(serverType):
		cc = server.NewServerComponent(name, spv.Interface().(server.Server))
	case spt.Implements(clientType):
		cc = client.NewClientComponent(name, spv.Interface().(client.Client))
	default:
		return nil, errors.New("set option should be a server or a client")
	}

	var st = spt.Elem()
	var sv = spv.Elem()

	for i := range st.NumField() {
		var sf = st.Field(i)
		if !sf.IsExported() {
			continue
		}

		var ft = sf.Type
		if ft.Kind() != reflect.Pointer {
			ft = reflect.PointerTo(ft)
		}

		var componentName string
		switch {
		case ft.Implements(serverType):
			var tags = strings.Split(sf.Tag.Get("server"), ",")
			componentName = tags[0]
		case ft.Implements(clientType):
			var tags = strings.Split(sf.Tag.Get("client"), ",")
			componentName = tags[0]
		default:
			continue
		}

		if componentName == "" {
			componentName = util.LowerCamelCase(sf.Name)
		}

		var val = option.Option(opt).Get(componentName)
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

		var scs []component.Component
		scs, err = (Option(tmp)).reflectSetOption(componentName, fv, ft)
		if err != nil {
			return nil, err
		}

		cs = append(cs, scs...)
	}

	st = spt
	sv = spv

	sm, ok := st.MethodByName("SetOption")
	if !ok {
		st = st.Elem()
		sv = sv.Elem()
		sm, ok = st.MethodByName("SetOption")
	}

	if ok {
		var mt = sm.Type
		if mt.NumIn() != 2 {
			return nil, errors.Errorf("invalid parameter of the server '%s' option set funcation", st.Name())
		}

		var at = mt.In(1)
		if at.Kind() != reflect.Struct && at.Elem().Kind() != reflect.Struct {
			return nil, errors.Errorf("the parameter of the server '%s' option set funcation should be a struct or struct pointer", st.Name())
		}

		if at.Kind() == reflect.Pointer {
			var av = reflect.New(at.Elem())
			err = jsonx.Object(opt).ToStruct(av.Interface())
			if err != nil {
				return nil, err
			}

			sv.Method(sm.Index).Call([]reflect.Value{av})
		} else {
			var av = reflect.New(at)
			err = jsonx.Object(opt).ToStruct(av.Interface())
			if err != nil {
				return nil, err
			}
			sv.Method(sm.Index).Call([]reflect.Value{av.Elem()})
		}
	}

	cs = append(cs, cc)
	return cs, nil
}
