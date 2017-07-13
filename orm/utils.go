package orm

import (
	// `reflect`
	"strings"

	"github.com/xslower/goutils/utils"
)

// convert camel string to underlined, XxYy to xx_yy
func ToUnderline(cameled string) string {
	underlined := strings.ToLower(string(cameled[0]))
	for _, c := range cameled[1:] {
		if isUpperCase(c) {
			underlined += `_`
			c += 32
		}
		underlined += string(c)
	}
	return underlined
}

func isUpperCase(c rune) bool {
	if c >= 'A' && c <= 'Z' {
		return true
	}
	return false
}

// convert underlined string to camel, aa_bb to AaBb
func ToCamel(underlined string) string {
	slice := strings.Split(underlined, `_`)
	if len(slice) == 1 { //if have no underline, then return
		return underlined
	}
	cameled := ``
	for _, elem := range slice {
		d := elem[0]
		if d >= 'a' && d <= 'z' {
			d -= 32
		}
		cameled += string(d) + elem[1:]
	}
	return cameled
}

//not support int8/uint8 please use int16/uint16
func InterfaceToString(ifc interface{}) (value string, stype uint8) {
	return utils.InterfaceToString(ifc)

}

//not support int8/uint8 please use int16/uint16
func InterfaceToInt(ifc interface{}) (value int, stype uint8) {
	return utils.InterfaceToInt(ifc)

}
