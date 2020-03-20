package utils

import "strings"

func RemoveQuotationMarks(str []byte) []byte {
	strLen := len(str)
	if strLen >= 2 && strings.HasPrefix(string(str), "\"") &&
		strings.HasSuffix(string(str), "\"") {
		str = str[1 : strLen-1]
	}
	return str
}
