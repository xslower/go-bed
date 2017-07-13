package orm

import (
	// `fmt`
	// `strconv`
	"strings"
)

const (
	AND = `AND`
	OR  = `OR`
)

type ICondition interface {
	GetSql() string
	Columns() []string
	Values() []interface{}
	Get(key string) (value interface{}, exist bool)
	Del(key string) bool
}

func NewConds() *Condition {
	c := &Condition{rel: AND, prepare: true}
	return c
}

func NewCondsOr() *Condition {
	c := &Condition{rel: OR, prepare: true}
	return c
}

type Condition struct {
	rel     string
	exprs   []ISubExpr
	conds   []ICondition
	prepare bool
}

func (this *Condition) InitWithMap(condMap map[string]interface{}) *Condition {
	for key, val := range condMap {
		this.exprs = append(this.exprs, &ExprEqual{key, val})
	}
	return this
}

func (this *Condition) InitWithSlice(columns []string, values []interface{}) *Condition {
	if len(columns) == len(values) {
		for i, clm := range columns {
			this.exprs = append(this.exprs, &ExprEqual{clm, values[i]})
		}
	} else {
		panic(`given values number is not match columns number`)
	}
	return this
}

func (this *Condition) AddConds(c ICondition) *Condition {
	this.conds = append(this.conds, c)
	return this
}

func (this *Condition) Equal(column string, value interface{}) *Condition {
	this.exprs = append(this.exprs, &ExprEqual{column, value})
	return this
}

func (this *Condition) Between(column string, start, end interface{}) *Condition {
	this.exprs = append(this.exprs, &ExprBetween{column, start, end})
	return this
}

func (this *Condition) In(column string, elements ...interface{}) *Condition {
	this.exprs = append(this.exprs, &ExprIn{column, elements})
	return this
}

func (this *Condition) Like(column string, value string) *Condition {
	this.exprs = append(this.exprs, &ExprLike{column, value})
	return this
}

func (this *Condition) Op(column string, value interface{}, op string) *Condition {
	this.exprs = append(this.exprs, &ExprOp{column, op, value})
	return this
}

func (this *Condition) GetSql() string {
	exprSlice := []string{}
	for _, expr := range this.exprs {
		exprSlice = append(exprSlice, expr.Sql(this.prepare))
	}
	for _, cond := range this.conds {
		exprSlice = append(exprSlice, `(`+cond.GetSql()+`)`)
	}
	sql := strings.Join(exprSlice, ` `+this.rel+` `)
	return sql
}

func (this *Condition) Values() []interface{} {
	ifcs := []interface{}{}
	for _, expr := range this.exprs {
		ifcs = append(ifcs, expr.Values()...)
	}
	for _, cond := range this.conds {
		ifcs = append(ifcs, cond.Values()...)
	}
	return ifcs
}

func (this *Condition) Columns() []string {
	clmSlice := []string{}
	for _, expr := range this.exprs {
		clmSlice = append(clmSlice, expr.Column())
	}
	for _, cond := range this.conds {
		clmSlice = append(clmSlice, cond.Columns()...)
	}
	return clmSlice
}

func (this *Condition) Get(key string) (value interface{}, exist bool) {
	for _, expr := range this.exprs {
		if expr.Column() == key {
			// vals :=
			return expr.Values()[0], true
		}
	}
	for _, cond := range this.conds {
		value, exist = cond.Get(key)
		if exist {
			return
		}
	}
	return nil, false
}

func (this *Condition) Del(key string) bool {
	for i, expr := range this.exprs {
		if expr.Column() == key {
			this.exprs = append(this.exprs[:i], this.exprs[i+1:]...)
			return true
		}
	}
	for _, cond := range this.conds {
		if cond.Del(key) {
			return true
		}
	}
	return false
}

type ISubExpr interface {
	Sql(prepare bool) string
	Values() []interface{}
	Column() string
}

type ExprEqual struct {
	column string
	value  interface{}
}

func (this *ExprEqual) Sql(prepare bool) string {
	val := `?`
	if !prepare {
		val = IfcToSqlValue(this.value)
	}
	expr := "`" + this.column + "` = " + val
	return expr
}

func (this *ExprEqual) Values() []interface{} {
	return []interface{}{this.value}
}

func (this *ExprEqual) Column() string {
	return this.column
}

type ExprIn struct {
	column string
	values []interface{}
}

func (this *ExprIn) Sql(prepare bool) string {
	ss := []string{}
	for _, val := range this.values {
		v := `?`
		if !prepare {
			v = IfcToSqlValue(val)
		}
		ss = append(ss, v)
	}
	expr := "`" + this.column + "`" + ` IN (` + strings.Join(ss, `, `) + `)`
	return expr
}

func (this *ExprIn) Values() []interface{} {
	return this.values
}

func (this *ExprIn) Column() string {
	return this.column
}

type ExprBetween struct {
	column string
	start  interface{}
	end    interface{}
}

func (this *ExprBetween) Sql(prepare bool) string {
	s := `?`
	e := `?`
	if !prepare {
		s = IfcToSqlValue(this.start)
		e = IfcToSqlValue(this.end)
	}
	expr := "`" + this.column + "`" + ` BETWEEN ` + s + ` AND ` + e
	return expr
}

func (this *ExprBetween) Values() []interface{} {
	return []interface{}{this.start, this.end}
}

func (this *ExprBetween) Column() string {
	return this.column
}

type ExprLike struct {
	column string
	value  string
}

func (this *ExprLike) Sql(prepare bool) string {
	val := `?`
	if !prepare {
		val = this.value
		if strings.Index(val, `%`) < 0 {
			val = `%` + val + `%`
		}
	}
	expr := "`" + this.column + "`" + ` LIKE ` + val
	return expr
}

func (this *ExprLike) Values() []interface{} {
	return []interface{}{this.value}
}

func (this *ExprLike) Column() string {
	return this.column
}

var ValidOp = map[string]bool{`>`: true, `>=`: true, `<`: true, `<=`: true}

type ExprOp struct {
	column string
	op     string
	value  interface{}
}

func (this *ExprOp) Sql(prepare bool) string {
	val := `?`
	if !prepare {
		val = IfcToSqlValue(this.value)
	}
	expr := "`" + this.column + "` " + this.op + ` ` + val
	return expr
}

func (this *ExprOp) Values() []interface{} {
	return []interface{}{this.value}
}

func (this *ExprOp) Column() string {
	return this.column
}
