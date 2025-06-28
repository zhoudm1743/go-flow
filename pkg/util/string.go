package util

import (
	"bytes"
	"strings"
	"unicode"
)

var StringUtil = stringUtil{}

// arrayUtil 数组工具类
type stringUtil struct{}

func (su stringUtil) ToSnakeCase(s string) string {
	buf := bytes.Buffer{}
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				buf.WriteRune('_')
			}
			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}
func (su stringUtil) ToCamelCase(s string) string {
	words := strings.Split(s, "_")
	for i := 1; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

func (su stringUtil) HexStringToByte(hexStr string) byte {
	result := byte(0)
	for i, c := range hexStr {
		asciiVal := byte(c)
		var hexVal byte
		if asciiVal >= '0' && asciiVal <= '9' {
			hexVal = asciiVal - '0'
		} else if asciiVal >= 'A' && asciiVal <= 'F' {
			hexVal = asciiVal - 'A' + 10
		} else if asciiVal >= 'a' && asciiVal <= 'f' {
			hexVal = asciiVal - 'a' + 10
		}
		shiftAmount := (len(hexStr) - 1 - i) * 4
		result |= (hexVal << shiftAmount)
	}
	return result
}
