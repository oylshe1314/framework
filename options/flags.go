package options

import (
	"github.com/oylshe1314/framework/errors"
	"strings"
)

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
