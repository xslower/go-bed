package orm

import (
	// `reflect`
	`strconv`
	`strings`
	`time`
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

const (
	TYPE_STRING uint8 = iota
	TYPE_INT
	TYPE_UINT
	TYPE_FLOAT
	TYPE_BOOL
	TYPE_NULL
)

//not support int8/uint8 please use int16/uint16
func InterfaceToString(ifc interface{}) (value string, stype uint8) {
	if ifc == nil {
		return ``, TYPE_NULL
	}
	switch val := ifc.(type) {
	case int16:
		return strconv.FormatInt(int64(val), 10), TYPE_INT
	case int32:
		return strconv.FormatInt(int64(val), 10), TYPE_INT
	case int64:
		return strconv.FormatInt(int64(val), 10), TYPE_INT
	case int:
		return strconv.FormatInt(int64(val), 10), TYPE_INT
	case uint16:
		return strconv.FormatUint(uint64(val), 10), TYPE_UINT
	case uint32:
		return strconv.FormatUint(uint64(val), 10), TYPE_UINT
	case uint64:
		return strconv.FormatUint(uint64(val), 10), TYPE_UINT
	case uint:
		return strconv.FormatUint(uint64(val), 10), TYPE_UINT
	case float32:
		return strconv.FormatFloat(float64(val), 'g', -1, 32), TYPE_FLOAT
	case float64:
		return strconv.FormatFloat(val, 'g', -1, 64), TYPE_FLOAT
	case bool:
		ret := `0`
		if val {
			ret = `1`
		}
		return ret, TYPE_BOOL
	case string:
		return val, TYPE_STRING
	case []byte:
		return string(val), TYPE_STRING
	case time.Time:
		return ``, TYPE_STRING
	default:
		panic(`Not support field type`)
	}

}

//not support int8/uint8 please use int16/uint16
func InterfaceToInt(ifc interface{}) (value int, stype uint8) {
	if ifc == nil {
		return 0, TYPE_NULL
	}
	switch val := ifc.(type) {
	case int16:
		return int(val), TYPE_INT
	case int32:
		return int(val), TYPE_INT
	case int64:
		return int(val), TYPE_INT
	case int:
		return val, TYPE_INT
	case uint16:
		return int(val), TYPE_UINT
	case uint32:
		return int(val), TYPE_UINT
	case uint64:
		return int(val), TYPE_UINT
	case uint:
		return int(val), TYPE_UINT
	case float32:
		return int(val), TYPE_FLOAT
	case float64:
		return int(val), TYPE_FLOAT
	case bool:
		ret := 0
		if val {
			ret = 1
		}
		return ret, TYPE_BOOL
	case string:
		ret, _ := strconv.Atoi(val)
		return ret, TYPE_STRING
	case []byte:
		ret, _ := strconv.Atoi(string(val))
		return ret, TYPE_STRING
	case time.Time:
		return 0, TYPE_STRING
	default:
		panic(`Not support field type`)
	}

}

func IfcToSqlValue(ifc interface{}) string {
	val, typ := InterfaceToString(ifc)
	if typ == TYPE_STRING {
		return `'` + val + `'`
	}
	return val
}
