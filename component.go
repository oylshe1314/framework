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

func ScanComponents(com any, opt option.Option) ([]component.Component, error) {
	var cpt = reflect.TypeOf(com)
	if cpt.Kind() != reflect.Pointer || cpt.Elem().Kind() != reflect.Struct {
		return nil, errors.New("server should be a struct pointer")
	}
	return scanComponents(util.LowerCamelCase(cpt.Elem().Name()), reflect.ValueOf(com), cpt, opt)
}

func scanComponents(name string, cpv reflect.Value, cpt reflect.Type, opt option.Option) ([]component.Component, error) {
	var err error
	var cs []component.Component

	var serverType = reflect.TypeFor[*server.Server]().Elem()
	var clientType = reflect.TypeFor[*client.Client]().Elem()

	var cc component.Component
	switch {
	case cpt.Implements(serverType):
		cc = server.NewServerComponent(name, cpv.Interface().(server.Server))
	case cpt.Implements(clientType):
		cc = client.NewClientComponent(name, cpv.Interface().(client.Client))
	default:
		return nil, errors.New("set option should be a server or a client")
	}

	if opt != nil {
		som, ok := cpt.MethodByName("SetOption")
		if ok {
			var mt = som.Type
			if mt.NumIn() != 2 {
				return nil, errors.Errorf("invalid parameter of the server '%s' option set function", cpt.String())
			}

			var at = mt.In(1)
			if at.Kind() != reflect.Struct && at.Elem().Kind() != reflect.Struct {
				return nil, errors.Errorf("the parameter of the server '%s' option set function should be a struct or struct pointer", cpt.String())
			}

			if at.Kind() == reflect.Pointer {
				var av = reflect.New(at.Elem())
				err = jsonx.Object(opt).ToStruct(av.Interface())
				if err != nil {
					return nil, err
				}

				cpv.Method(som.Index).Call([]reflect.Value{av})
			} else {
				var av = reflect.New(at)
				err = jsonx.Object(opt).ToStruct(av.Interface())
				if err != nil {
					return nil, err
				}
				cpv.Method(som.Index).Call([]reflect.Value{av.Elem()})
			}
		}
	}

	var ct = cpt.Elem()
	var cv = cpv.Elem()
	for i := range ct.NumField() {
		var sf = ct.Field(i)
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

		var subOpt = option.Option(nil)
		if opt != nil {
			var val = opt.Get(componentName)
			if val != nil {
				tmp, ok := val.(map[string]any)
				if !ok {
					return nil, errors.Errorf("can not set option '%s' of %s to server '%s' ", name, reflect.TypeOf(tmp).String(), sf.Name)
				}
				subOpt = tmp
			}
		}

		var fv = cv.Field(i)
		if sf.Type.Kind() == reflect.Pointer {
			if fv.IsNil() {
				return nil, errors.Errorf("the pointer of the server '%s' is nil", sf.Name)
			}
		} else {
			fv = fv.Addr()
		}

		var scs []component.Component
		scs, err = scanComponents(componentName, fv, ft, subOpt)
		if err != nil {
			return nil, err
		}

		cs = append(cs, scs...)
	}

	cs = append(cs, cc)
	return cs, nil
}
