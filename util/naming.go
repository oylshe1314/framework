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
	var ss []string
	var rs = []rune(s)

	var word []rune
	var uppers int
	var lowers int
	for _, r := range rs {
		if unicode.IsUpper(r) {
			if lowers > 0 {
				ss = append(ss, string(word))
				word = []rune{r}
				lowers = 0
				uppers = 1
			} else {
				word = append(word, r)
				uppers += 1
			}
		} else {
			if uppers > 1 {
				ss = append(ss, string(word[:uppers-1]))
				word = word[uppers-1:]
				uppers = 1
			}
			word = append(word, r)
			lowers += 1
		}
	}
	if len(word) > 0 {
		ss = append(ss, string(word))
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
