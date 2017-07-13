package orm

import (
	"strconv"
	"strings"
)

const (
	DEFAULT direct = 0
	ASC     direct = 1
	DESC    direct = -1
)

var (
	directMap = map[direct]string{1: `ASC`, -1: `DESC`}
)

type direct int8

func (this direct) String() string {
	return directMap[this]
}

type IAppends interface {
	GetSql() string
	Values() []interface{}
}

type ApdGroup struct {
	fields []string
}

func (this *ApdGroup) Add(fields ...string) {
	this.fields = append(this.fields, fields...)
}
func (this *ApdGroup) GetSql() string {
	if len(this.fields) == 0 {
		return ``
	}
	return ` GROUP BY ` + "`" + strings.Join(this.fields, "`, `") + "`"
}
func (this *ApdGroup) Values() []interface{} {
	return []interface{}{}
}

type ApdOrder struct {
	columns []string
	ds      []direct
}

func (this *ApdOrder) Add(field string, d direct) {
	if d == 0 {
		d = ASC
	}
	this.columns = append(this.columns, field)
	this.ds = append(this.ds, d)
}
func (this *ApdOrder) GetSql() string {
	if len(this.columns) == 0 {
		return ``
	}
	sql := ` ORDER BY `
	for i, clm := range this.columns {
		sql += "`" + clm + "` " + this.ds[i].String() + `, `
	}
	return sql[:len(sql)-2]
}
func (this *ApdOrder) Values() []interface{} {
	return []interface{}{}
}

type ApdLimit struct {
	offset  int
	num     int
	prepare bool
}

func (this *ApdLimit) GetSql() string {
	if this.offset == 0 && this.num == 0 {
		return ``
	}
	ofs := `?`
	nm := `?`
	if !this.prepare {
		ofs = strconv.Itoa(this.offset)
		nm = strconv.Itoa(this.num)
	}
	return ` LIMIT ` + ofs + `, ` + nm
}

func (this *ApdLimit) Values() []interface{} {
	if this.offset == 0 && this.num == 0 {
		return []interface{}{}
	}
	return []interface{}{this.offset, this.num}
}

func NewAppends() *Appends {
	a := &Appends{prepare: true}
	a.group = &ApdGroup{}
	a.order = &ApdOrder{}
	a.limit = &ApdLimit{prepare: a.prepare}
	return a
}

type Appends struct {
	group   *ApdGroup
	order   *ApdOrder
	limit   *ApdLimit
	prepare bool
}

func (this *Appends) GroupBy(fields ...string) *Appends {
	this.group.Add(fields...)
	return this
}
func (this *Appends) OrderBy(field string, d direct) *Appends {
	this.order.Add(field, d)
	return this
}
func (this *Appends) Limit(offset, num int) *Appends {
	this.limit.offset = offset
	this.limit.num = num

	return this
}
func (this *Appends) GetSql() string {
	return this.group.GetSql() + this.order.GetSql() + this.limit.GetSql()
}
func (this *Appends) Values() []interface{} {
	return this.limit.Values()
}
