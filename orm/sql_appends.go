package orm

import (
	`strconv`
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
	return ` GROUP BY ` + "`" + strings.Join(this.fields, "`, `") + "`"
}
func (this *ApdGroup) Values() []interface{} {
	return []interface{}{}
}

type ApdOrder struct {
	fields map[string]direct
}

func (this *ApdOrder) Add(field string, d direct) {
	if d == 0 {
		d = ASC
	}
	this.fields[field] = d
}
func (this *ApdOrder) GetSql() string {
	sql := ` ORDER BY `
	for key, val := range this.fields {
		sql += "`" + key + "` " + val.String() + `, `
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
	ofs := `?`
	nm := `?`
	if !this.prepare {
		ofs = strconv.Itoa(this.offset)
		nm = strconv.Itoa(this.num)
	}
	return ` LIMIT ` + ofs + `, ` + nm
}

func NewAppends() *Appends {
	return &Appends{prepare: true}
}

type Appends struct {
	group   *ApdGroup
	order   *ApdOrder
	limit   *ApdLimit
	prepare bool
}

func (this *Appends) GroupBy(fields ...string) *Appends {
	if this.group == nil {
		this.group = &ApdGroup{fields}
	} else {
		this.group.Add(fields...)
	}
	return this
}
func (this *Appends) OrderBy(field string, d direct) *Appends {
	if this.order == nil {
		this.order = &ApdOrder{}
	}
	this.order.Add(field, d)
	return this
}
func (this *Appends) Limit(offset, num int) *Appends {
	if this.limit == nil {
		this.limit = &ApdLimit{offset, num, this.prepare}
	} else {
		this.limit.offset = offset
		this.limit.num = num
	}
	return this
}
func (this *Appends) GetSql() string {
	return this.group.GetSql() + this.order.GetSql() + this.limit.GetSql()
}
func (this *Appends) Values() []interface{} {
	return []interface{}{this.limit.offset, this.limit.num}
}
