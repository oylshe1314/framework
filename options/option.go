package options

import (
	"fmt"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/log"
	"github.com/oylshe1314/framework/util"
	"reflect"
	"sort"
	"strings"
)

type Optional interface {
	Init() error
}

func InitServers(ss ...Optional) (err error) {
	for _, s := range ss {
		err = s.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

type Options map[string]interface{}

func (options Options) Get(name string) interface{} {
	return options[name]
}

func (options Options) Put(name string, value interface{}) {
	options[name] = value
}

func (options Options) Merge(newOptions Options) {
	if newOptions == nil {
		return
	}
	for name, value := range newOptions {
		ov, ok1 := options[name]
		if ok1 {
			no, ok2 := value.(map[string]interface{})
			if ok2 {
				oo, ok3 := ov.(map[string]interface{})
				if ok3 {
					Options(oo).Merge(no)
					continue
				}
			}
		}
		options[name] = value
	}
}

func (options Options) setOptions(osv reflect.Value, parent string) error {
	var ost = osv.Type()
	var mn = ost.NumMethod()
	for mi := 0; mi < mn; mi++ {
		var m = ost.Method(mi)
		if !m.IsExported() {
			continue
		}

		if strings.HasPrefix(m.Name, "With") {
			var name = m.Name[4:]
			var key = util.LowerCamelCase(name)
			var val = options.Get(key)
			if val == nil {
				continue
			}

			if m.Type.NumIn() != 2 {
				continue
			}

			var pt = m.Type.In(1)
			var pv, err = util.NewReflectValueFromJson(val, pt)
			if err != nil {
				return errors.Errorf("set options '%s' failed, %v", parent+name, err)
			}

			m.Func.Call([]reflect.Value{osv, pv})
		}
	}

	if ost.Kind() == reflect.Pointer {
		ost = ost.Elem()
		osv = osv.Elem()
	}

	if ost.Kind() != reflect.Struct {
		return nil
	}

	var fn = ost.NumField()
	for fi := 0; fi < fn; fi++ {
		var f = ost.Field(fi)
		if !f.IsExported() {
			continue
		}

		if strings.HasSuffix(f.Name, "Server") || strings.HasSuffix(f.Name, "Client") || strings.HasSuffix(f.Name, "Config") {
			var name = util.LowerCamelCase(f.Name)
			var val = options.Get(name)
			if val == nil {
				continue
			}

			subOptions, ok := val.(map[string]interface{})
			if !ok {
				continue
			}

			fv := osv.Field(fi)
			if fv.Kind() != reflect.Pointer {
				fv = fv.Addr()
			}

			if _, ok = fv.Interface().(Optional); !ok {
				continue
			}

			var err = Options(subOptions).setOptions(fv, parent+name+".")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (options Options) SetOptions(server Optional) error {
	return options.setOptions(reflect.ValueOf(server), "")
}

func (options Options) Init(svr Optional) (err error) {
	err = options.SetOptions(svr)
	if err != nil {
		return err
	}

	return svr.Init()
}
func LogOptions(logger log.Logger, options Options) {
	var list util.Pairs[*[]string, interface{}]
	collectOptions(logger, nil, options, &list)
	sort.Slice(list, func(i, j int) bool { //why sort here, restless for overeating?
		var ki, kj = *list[i].Key, *list[j].Key
		if len(ki) == len(kj) {
			for kii := range ki {
				if ki[kii] != kj[kii] {
					return ki[kii] < kj[kii]
				}
			}
			return true
		} else {
			if len(ki) == 1 {
				return true
			}

			if len(kj) == 1 {
				return false
			}

			if ki[0] == kj[0] {
				return len(ki) < len(kj)
			} else {
				return ki[0] < kj[0]
			}
		}
	})
	for _, pair := range list {
		logger.Info("Option: ", strings.Join(*pair.Key, "."), " = ", pair.Value)
	}
}

func collectOptions(logger log.Logger, p []string, options Options, list *util.Pairs[*[]string, interface{}]) {
	for k, v := range options {
		var np = make([]string, len(p)+1)
		copy(np, p)
		np[len(p)] = k
		switch rv := v.(type) {
		case map[string]interface{}:
			collectOptions(logger, np, rv, list)
		case []interface{}:
			collectOptionArray(logger, np, rv, list)
		default:
			*list = append(*list, util.NewPair(&np, v))
		}
	}
}

func collectOptionArray(logger log.Logger, p []string, array []interface{}, list *util.Pairs[*[]string, interface{}]) {
	for i, v := range array {
		var np = make([]string, len(p))
		copy(np, p[:len(p)-1])
		np[len(p)-1] = fmt.Sprintf("%s[%d]", p[len(p)-1], i)
		switch rv := v.(type) {
		case map[string]interface{}:
			collectOptions(logger, np, rv, list)
		case []interface{}:
			collectOptionArray(logger, np, rv, list)
		default:
			*list = append(*list, util.NewPair(&np, rv))
		}
	}
}
