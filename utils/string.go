package utils

import (
	"strings"
	"unsafe"
)

func ItHas(str, cnt string) bool {
	return strings.Contains(str, cnt)
}

func ItHasSuffix(str, cnt string) bool {
	return strings.HasSuffix(str, cnt)
}

func SplitString(s, sep string) (ss []string) {
	if s = strings.TrimSpace(s); s == "" {
		return
	}

	for _, val := range strings.Split(s, sep) {
		if val = strings.TrimSpace(val); val != "" {
			ss = append(ss, val)
		}
	}
	return
}

func Join(sources []string, delimeter, defaultValue string) string {
	if len(sources) == 0 {
		return defaultValue
	}

	return strings.Join(sources, delimeter)
}

func ByteToString(input []byte) string {
	return *(*string)(unsafe.Pointer(&input))
}
