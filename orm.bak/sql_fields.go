package orm

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
	fi.values = [][]interface{}{}
	// fi.values[0] = make([]interface{}, 1)
	return fi
}
func NewFieldsMultiInsertWithMap(fields map[string]interface{}) *FieldsMultiInsert {
	fi := NewFieldsMultiInsert()
	value := []interface{}{}
	for clm, val := range fields {
		fi.columns = append(fi.columns, clm)
		value = append(value, val)
	}
	fi.values = append(fi.values, value)
	return fi
}
func NewFieldsMultiInsertWithSlice(columns []string, values []interface{}) *FieldsMultiInsert {
	fi := NewFieldsMultiInsert()
	if len(columns) != len(values) {
		panic(`given values number is not match columns number`)
	}
	fi.columns = columns
	fi.values = append(fi.values, values)
	return fi
}

type FieldsMultiInsert struct {
	columns []string
	values  [][]interface{} //用list目的是可以支持多条插入。
	prepare bool
}

//增加一条记录
func (this *FieldsMultiInsert) AddMap(fields map[string]interface{}) *FieldsMultiInsert {
	num := len(this.columns)
	if num == 0 {
		num = len(fields)
	} else if num != len(fields) {
		panic(`given fields number is not match columns number`)

	}
	values := make([]interface{}, num)
	for i, clm := range this.columns {
		values[i] = fields[clm]
	}
	this.values = append(this.values, values)
	return this
}

//增加一条记录，
func (this *FieldsMultiInsert) AddSlice(values []interface{}) *FieldsMultiInsert {
	if len(values) != len(this.columns) {
		panic(`given values number is not match columns number`)
	}
	this.values = append(this.values, values)
	return this
}
func (this *FieldsMultiInsert) GetSql() string {
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
	ret := this.values[0]
	for _, val := range this.values[1:] {
		ret = append(ret, val...)
	}
	return ret
}
func (this *FieldsMultiInsert) Columns() []string {
	return this.columns
}
func (this *FieldsMultiInsert) Get(key string) (value interface{}, exist bool) {
	for i, clm := range this.columns {
		if clm == key {
			return this.values[0][i], true
		}
	}
	return nil, false
}
func (this *FieldsMultiInsert) Del(key string) bool {
	n := -1
	for i, clm := range this.columns {
		if clm == key {
			this.columns = append(this.columns[:i], this.columns[i+1:]...)
			n = i
		}
	}
	if n == -1 {
		return false
	}
	for i, val := range this.values {
		this.values[i] = append(val[:n], val[n+1:]...)
	}
	return true
}

/**
 * FieldWrite
 */

func NewFieldsWrite() *FieldsWrite {
	return &FieldsWrite{prepare: true}
}
func NewFieldsWriteWithMap(fields map[string]interface{}) *FieldsWrite {
	fu := NewFieldsWrite()
	for key, val := range fields {
		fu.fldList = append(fu.fldList, &Field{key, val, _VALTYPE_NORMAL})
	}
	return fu
}
func NewFieldsWriteWithSlice(columns []string, values []interface{}) *FieldsWrite {
	fu := NewFieldsWrite()
	if len(columns) == len(values) {
		for i, clm := range columns {
			fu.fldList = append(fu.fldList, &Field{clm, values[i], _VALTYPE_NORMAL})
		}
	} else {
		panic(`given values number is not match columns number`)
	}
	return fu
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
