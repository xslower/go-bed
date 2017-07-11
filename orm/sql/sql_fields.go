package sql

import (
	"strings"
)

type pSqlType uint8

const (
	_VALTYPE_NORMAL pSqlType = iota
	_VALTYPE_OP_SELF_ADD
	_VALTYPE_OP_SELF_SUB
)

type IFields interface {
	GetSql() string
	Columns() []string
	Values() []interface{}
	Get(key string) (value interface{}, exist bool)
	Del(key string) bool
}

/**
* MultiInsert
 */
func NewFieldsMultiInsert() *FieldsMultiInsert {
	fi := &FieldsMultiInsert{prepare: true}
	fi.colIndex = map[string]int{}
	fi.columns = []string{}
	fi.values = [][]interface{}{}

	// fi.values[0] = make([]interface{}, 1)
	return fi
}

type FieldsMultiInsert struct {
	colIndex map[string]int
	columns  []string
	values   [][]interface{} //用list目的是可以支持多条插入。
	prepare  bool
}

func (this *FieldsMultiInsert) InitWithMap(fields map[string]interface{}) *FieldsMultiInsert {
	value := []interface{}{}
	for clm, val := range fields {
		this.columns = append(this.columns, clm)
		value = append(value, val)
	}
	this.values = append(this.values, value)
	for i, clm := range this.columns {
		this.colIndex[clm] = i
	}
	return this
}

func (this *FieldsMultiInsert) InitWithSlice(columns []string, values []interface{}) *FieldsMultiInsert {
	if len(columns) != len(values) {
		panic(`given values number is not match columns number`)
	}
	this.columns = columns
	this.values = append(this.values, values)
	for i, clm := range this.columns {
		this.colIndex[clm] = i
	}
	return this
}

//增加一条记录。需要处理fields中字段数量跟现有字段数量不匹配的情况。
func (this *FieldsMultiInsert) AddMap(fields map[string]interface{}) *FieldsMultiInsert {
	num := len(this.columns)
	if num == 0 {
		return this.InitWithMap(fields)
	}
	for clm, _ := range fields {
		_, ok := this.colIndex[clm]
		if !ok {
			this.addColumn(clm)
			num++
		}
	}
	values := make([]interface{}, num)
	for i, clm := range this.columns {
		values[i] = fields[clm]
	}
	this.values = append(this.values, values)
	return this
}

//增加一条记录，
func (this *FieldsMultiInsert) AddSlice(values []interface{}, columns ...string) *FieldsMultiInsert {
	num := len(this.columns)
	if num == 0 {
		return this.InitWithSlice(columns, values)
	}
	if len(columns) == 0 {
		if len(values) != num {
			panic(`given [values]'s number is not match this.columns' number`)
		}
		this.values = append(this.values, values)
		return this
	}
	if len(values) != len(columns) {
		panic(`given [values]'s number is not match given [columns]'s number`)
	}
	for _, clm := range columns {
		_, ok := this.colIndex[clm]
		if !ok {
			this.addColumn(clm)
			num++
		}
	}
	val := make([]interface{}, num)
	for i, clm := range columns {
		val[this.colIndex[clm]] = values[i]
	}
	this.values = append(this.values, val)
	return this
}
func (this *FieldsMultiInsert) GetSql() string {
	if len(this.columns) == 0 {
		return ``
	}
	sql := `(` + strings.Join(this.columns, `, `) + `) VALUES `
	if this.prepare {
		quest := make([]string, len(this.columns))
		for i, _ := range quest {
			quest[i] = `?`
		}
		questStr := `(` + strings.Join(quest, `, `) + `)`
		for i := 0; i < len(this.values); i++ {
			sql += questStr + `, `
		}
		return sql[:len(sql)-2]
	}
	strRowSlice := make([]string, len(this.columns))
	for _, val := range this.values {
		for i, v := range val {
			strRowSlice[i] = IfcToSqlValue(v)
		}
		sql += `(` + strings.Join(strRowSlice, `, `) + `), `
	}
	return sql[:len(sql)-2]
}
func (this *FieldsMultiInsert) Values() []interface{} {
	if len(this.values) == 1 {
		return this.values[0]
	}
	ret := []interface{}{}
	for _, val := range this.values {
		ret = append(ret, val...)
	}
	return ret
}
func (this *FieldsMultiInsert) Columns() []string {
	return this.columns
}
func (this *FieldsMultiInsert) Get(key string) (value interface{}, exist bool) {
	index, ok := this.colIndex[key]
	if !ok {
		return nil, false
	}
	return this.values[0][index], true
}
func (this *FieldsMultiInsert) Del(key string) bool {
	index, ok := this.colIndex[key]
	if !ok {
		return false
	}
	delete(this.colIndex, key)
	this.columns = append(this.columns[:index], this.columns[index+1:]...)
	for i, val := range this.values {
		this.values[i] = append(val[:index], val[index+1:]...)
	}
	return true
}
func (this *FieldsMultiInsert) addColumn(clm string) {
	this.columns = append(this.columns, clm)
	this.colIndex[clm] = len(this.columns) - 1
	for i, _ := range this.values {
		this.values[i] = append(this.values[i], ``)
	}
}

/**
 * FieldWrite
 */

func NewFieldsWrite() *FieldsWrite {
	return &FieldsWrite{prepare: true}
}

type Field struct {
	column  string
	value   interface{}
	valType pSqlType
}

type FieldsWrite struct {
	fldList []*Field
	prepare bool
}

func (this *FieldsWrite) InitWithMap(fields map[string]interface{}) *FieldsWrite {
	for key, val := range fields {
		this.fldList = append(this.fldList, &Field{key, val, _VALTYPE_NORMAL})
	}
	return this
}

func (this *FieldsWrite) InitWithSlice(columns []string, values []interface{}) *FieldsWrite {
	if len(columns) == len(values) {
		for i, clm := range columns {
			this.fldList = append(this.fldList, &Field{clm, values[i], _VALTYPE_NORMAL})
		}
	} else {
		panic(`given values number is not match columns number`)
	}
	return this
}

func (this *FieldsWrite) Add(key string, value interface{}) *FieldsWrite {
	this.fldList = append(this.fldList, &Field{key, value, _VALTYPE_NORMAL})
	return this
}
func (this *FieldsWrite) AddSpl(key string, value interface{}, typ pSqlType) *FieldsWrite {
	this.fldList = append(this.fldList, &Field{key, value, typ})
	return this
}
func (this *FieldsWrite) GetSql() string {
	sql := ``
	if len(this.fldList) == 0 {
		return sql
	}
	for _, fld := range this.fldList {
		clm := fld.column
		sql += "`" + clm + "` = "
		switch fld.valType {
		case _VALTYPE_OP_SELF_ADD:
			sql += "`" + clm + "` + "
		case _VALTYPE_OP_SELF_SUB:
			sql += "`" + clm + "` - "
		}
		if this.prepare {
			sql += `?`
		} else {
			sql += IfcToSqlValue(fld.value)
		}
		sql += `, `
	}
	return sql[:len(sql)-2]
}
func (this *FieldsWrite) Values() []interface{} {
	ifcs := []interface{}{}
	for _, fld := range this.fldList {
		ifcs = append(ifcs, fld.value)
	}
	return ifcs
}
func (this *FieldsWrite) Columns() []string {
	columns := []string{}
	for _, fld := range this.fldList {
		columns = append(columns, fld.column)
	}
	return columns
}
func (this *FieldsWrite) Get(key string) (value interface{}, exist bool) {
	for _, fld := range this.fldList {
		if fld.column == key {
			return fld.value, true
		}
	}
	return nil, false
}
func (this *FieldsWrite) Del(key string) bool {
	for i, fld := range this.fldList {
		if fld.column == key {
			this.fldList = append(this.fldList[:i], this.fldList[i+1:]...)
			return true
		}
	}
	return false
}

func NewFieldsRead(columns ...string) *FieldsRead {
	return &FieldsRead{columns}
}

type FieldsRead struct {
	columns []string
}

func (this *FieldsRead) Add(column string) *FieldsRead {
	this.columns = append(this.columns, column)
	return this
}
func (this *FieldsRead) GetSql() string {
	if len(this.columns) > 0 {
		return strings.Join(this.columns, `, `)
	} else {
		return `*`
	}

}
func (this *FieldsRead) Values() []interface{} {
	return []interface{}{}
}
func (this *FieldsRead) Columns() []string {
	return this.columns
}
func (this *FieldsRead) Get(key string) (value interface{}, exist bool) {
	for _, clm := range this.columns {
		if clm == key {
			return nil, true
		}
	}
	return nil, false
}
func (this *FieldsRead) Del(key string) bool {
	for i, clm := range this.columns {
		if clm == key {
			this.columns = append(this.columns[:i], this.columns[i+1:]...)
			return true
		}
	}
	return false
}
