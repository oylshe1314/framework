package util

import (
	"strings"
	"unicode"
)

func toCamelCase(s string, b int) string {
	var vs = strings.Split(s, "_")
	if len(vs) == 1 {
		var rs = []rune(vs[0])
		if b == 0 {
			rs[0] = unicode.ToLower(rs[0])
		} else {
			rs[0] = unicode.ToUpper(rs[0])
		}
		return string(rs)
	}

	var ss []string
	for i, v := range vs {
		if len(v) == 0 {
			continue
		}
		var rs = []rune(v)
		for ri := range rs {
			if i > b && ri == 0 {
				rs[ri] = unicode.ToUpper(rs[ri])
			} else {
				rs[ri] = unicode.ToLower(rs[ri])
			}
		}
		ss = append(ss, string(rs))
	}
	return strings.Join(ss, "")
}

// LowerCamelCase 转小驼峰
func LowerCamelCase(s string) string {
	return toCamelCase(s, 0)
}

// UpperCamelCase 转大驼峰
func UpperCamelCase(s string) string {
	return toCamelCase(s, -1)
}

// SplitCameCase 按驼峰切割
func SplitCameCase(s string) []string {
	var li = 0
	var ss []string
	var rs = []rune(s)
	for i, r := range rs {
		if unicode.IsUpper(r) && i-li > 0 {
			ss = append(ss, string(rs[li:i]))
			li = i
		}
	}

	if li < len(rs) {
		ss = append(ss, string(rs[li:]))
	}

	return ss
}

// LowerSnakeCase 转小蛇型
func LowerSnakeCase(s string) string {
	return strings.ToLower(strings.Join(SplitCameCase(s), "_"))
}

// UpperSnakeCase 转大蛇型
func UpperSnakeCase(s string) string {
	return strings.ToUpper(strings.Join(SplitCameCase(s), "_"))
}
