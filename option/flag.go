package option

import (
	"framework/errors"
	"net/url"
	"strings"
)

type FlagOption Option

func (option FlagOption) String() string {
	return Option(option).String()
}

func (option FlagOption) Set(arg string) error {
	var i = strings.IndexByte(arg, '=')
	if i <= 0 {
		return errors.Error("bad flag option syntax: ", arg)
	}

	var value = arg[i+1:]
	value, err := url.QueryUnescape(value)
	if err != nil {
		value = arg[i+1:]
	}

	option.setFieldChain(strings.Split(arg[:i], "."), 0, value)
	return nil
}

func (option FlagOption) setFieldChain(keysChain []string, i int, value string) {
	if i+1 == len(keysChain) {
		Option(option).Set(keysChain[i], value)
	} else {
		var subOptions map[string]any
		var v = option[keysChain[i]]
		if v == nil {
			subOptions = make(map[string]any)
			Option(option).Set(keysChain[i], subOptions)
		} else {
			var ok = false
			subOptions, ok = v.(map[string]any)
			if !ok {
				subOptions = make(map[string]any)
				Option(option).Set(keysChain[i], subOptions)
			}
		}
		FlagOption(subOptions).setFieldChain(keysChain, i+1, value)
	}
}
