package server

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/oylshe1314/framework/errors"
	"github.com/oylshe1314/framework/log"
	"github.com/oylshe1314/framework/util"
	"os"
	"reflect"
	"sort"
	"strings"
)

type OptionalServer interface {
	Init() error
}

func InitServers(svrs ...OptionalServer) (err error) {
	for _, svr := range svrs {
		err = svr.Init()
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

			if _, ok = fv.Interface().(OptionalServer); !ok {
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

func (options Options) SetOptions(server OptionalServer) error {
	return options.setOptions(reflect.ValueOf(server), "")
}

func (options Options) Init(svr OptionalServer) (err error) {
	err = options.SetOptions(svr)
	if err != nil {
		return err
	}

	return svr.Init()
}

func ReadOptions(filename string) (Options, error) {
	var options = map[string]interface{}{} //Remind me, don't use Options{} to init the variable
	if len(filename) == 0 {
		return options, nil
	}

	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf, &options)
	if err != nil {
		return nil, err
	}

	return options, nil
}

type FlagOptions []string

func (p *FlagOptions) String() string {
	return strings.Join(*p, "\n")
}

func (p *FlagOptions) Set(v string) error {
	*p = append(*p, v)
	return nil
}

func (p *FlagOptions) Get() any {
	return *p
}

func (p *FlagOptions) Parse() (Options, error) {
	var flagOptions = *p
	if len(flagOptions) == 0 {
		return nil, nil
	}

	var options = map[string]interface{}{}
	for _, flagOption := range flagOptions {
		var idx = strings.IndexByte(flagOption, '=')
		if idx <= 0 {
			return nil, errors.Error("bad flag options syntax: ", flagOption)
		}

		setFieldChain(strings.Split(flagOption[:idx], "."), 0, flagOption[idx+1:], options)
	}

	return options, nil
}

func setFieldChain(chains []string, idx int, value string, options Options) {
	if idx+1 == len(chains) {
		/**
		Do not overwrite the option value that has existed, we think
		the option value that was set first is the correct value you
		want to set, the second option value with the same name	may
		just be a small mistake you made.
		*/
		if _, ok := options[chains[idx]]; !ok {
			options[chains[idx]] = value
		}
	} else {
		var subOptions map[string]interface{}
		var v = options[chains[idx]]
		if v == nil {
			subOptions = map[string]interface{}{}
			options[chains[idx]] = subOptions
		} else {
			var ok = false
			subOptions, ok = v.(map[string]interface{})
			if !ok {
				return
			}
		}
		setFieldChain(chains, idx+1, value, subOptions)
	}
}

func logOptions(logger log.Logger, options Options) {
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
