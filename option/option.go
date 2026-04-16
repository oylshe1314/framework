package option

import (
	"fmt"
	"framework/store"
	"sort"
	"strings"
)

type Optional[Opt any] struct {
	opt *Opt
}

func (this *Optional[Opt]) SetOption(opt *Opt) {
	this.opt = opt
}

func (this *Optional[Opt]) GetOption() *Opt {
	return this.opt
}

type Option map[string]any

func (option Option) Get(key string) any {
	return option[key]
}

func (option Option) Set(key string, value any) {
	option[key] = value
}

// Merge
// It will overwrite the value that already exists.
// We consider the new value to be the final correct value you want.
// That's why you set it using the flag argument.
func (option Option) Merge(newOption Option) {
	if newOption == nil {
		return
	}
	for nn, nv := range newOption {
		ov, ok1 := option[nn]
		if ok1 {
			no, ok2 := nv.(map[string]any)
			if ok2 {
				oo, ok3 := ov.(map[string]any)
				if ok3 {
					Option(oo).Merge(no)
					continue
				}
			}
		}
		option[nn] = nv
	}
}

func (option Option) String() string {
	var pairs []*store.Pair[string, any]
	option.collectOption("", &pairs)

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Key < pairs[j].Key
	})

	var sb strings.Builder
	for _, pair := range pairs {
		sb.WriteString(pair.Key)
		sb.WriteString("=")
		sb.WriteString(fmt.Sprint(pair.Value))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (option Option) collectOption(p string, pairs *[]*store.Pair[string, any]) {
	for key, value := range option {
		if p != "" {
			key = p + "." + key
		}
		switch tv := value.(type) {
		case Option:
			tv.collectOption(key, pairs)
		case map[string]any:
			Option(tv).collectOption(key, pairs)
		case []any:
			option.collectArray(key, tv, pairs)
		default:
			*pairs = append(*pairs, store.NewPair(key, tv))
		}
	}
}

func (option Option) collectArray(p string, ary []any, pairs *[]*store.Pair[string, any]) {
	for i, value := range ary {
		var key = fmt.Sprintf("%s[%d]", p, i)
		switch tv := value.(type) {
		case Option:
			tv.collectOption(key, pairs)
		case map[string]any:
			Option(tv).collectOption(key, pairs)
		case []any:
			option.collectArray(key, tv, pairs)
		default:
			*pairs = append(*pairs, store.NewPair(key, tv))
		}
	}
}
