package parser

import "strings"

func GetParsedValue(s string) (string, bool) {
	if strings.HasPrefix(s, "+") {
		return "OK\n", true
	} else if strings.HasPrefix(s, "-ERR") {
		return strings.Split(s, "-ERR")[1][1:], false
	} else if strings.HasPrefix(s, "-nil") {
		return "nil\n", true
	} else if hasPrefixes(s, "$", ":", "*", "#") {
		return trim(s, 1), true
	} else {
		return s, true
	}
}

func trim(s string, i int) string {
	r := []rune(s)

	if len(r) >= i {
		return string(r[i:])
	}

	return s
}

func hasPrefixes(s string, args ...string) bool {
	var has bool

	for _, a := range args {
		has = strings.HasPrefix(s, a)
		if has {
			return true
		}
	}

	return has
}
