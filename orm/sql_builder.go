/**
* sql builder for mysql
 */
package orm

import (
// "strings"

)

var (
	fNewSqlBuilder   func(db, table string) ISqlBuilder
	gSqlBuFnRegistry = map[string]func(db, table string) ISqlBuilder{}
	// gSqlBuRegistry = map[string]ISqlBuilder
)

func init() {
	gSqlBuFnRegistry[`mysql`] = NewMysqlBuilder
}

type ISqlElem interface {
	GetSql() string
	Columns() []string
	Values() []interface{}
	Get(key string) (value interface{}, exist bool)
	Del(key string) bool
}

type ISqlBuilder interface {
	Select(fields IFields, condition ICondition, appends IAppends) string
	Insert(fields IFields) string
	InsertOrUpdate(fields IFields) string
	MultiInsert(fields IFields) string
	Update(fields IFields, condition ICondition) string
	Delete(condition ICondition) string
}

func NewMysqlBuilder(db, table string) ISqlBuilder {
	return &MysqlBuilder{db, table}
}

type MysqlBuilder struct {
	db    string
	table string
}

func (this *MysqlBuilder) Select(fields IFields, condition ICondition, appends IAppends) string {
	ftn := fullTableName(this.db, this.table)
	fldStr := `*`
	if fields != nil {
		fldStr = fields.GetSql()
	}
	condStr := ``
	if condition != nil {
		condStr = ` WHERE ` + condition.GetSql()
	}

	sql := `SELECT ` + fldStr + ` FROM ` + ftn + condStr
	return sql
}

func (this *MysqlBuilder) MultiInsert(fields IFields) string {
	ftn := fullTableName(this.db, this.table)
	sql := `INSERT ` + ftn + ` ` + fields.GetSql()
	return sql
}

func (this *MysqlBuilder) Insert(fields IFields) string {
	ftn := fullTableName(this.db, this.table)
	sql := `INSERT ` + ftn + ` SET ` + fields.GetSql()
	return sql
}

func (this *MysqlBuilder) InsertOrUpdate(fields IFields) string {
	ftn := fullTableName(this.db, this.table)
	fldStr := fields.GetSql()
	sql := `INSERT ` + ftn + ` SET ` + fldStr + ` ON DUPLICATE KEY UPDATE ` + fldStr
	return sql
}

func (this *MysqlBuilder) Update(fields IFields, condition ICondition) string {
	ftn := fullTableName(this.db, this.table)
	condStr := ``
	if condition != nil {
		condStr += ` WHERE ` + condition.GetSql()
	}
	sql := `UPDATE ` + ftn + ` SET ` + fields.GetSql() + condStr
	return sql
}

func (this *MysqlBuilder) Delete(condition ICondition) string {

	ftn := fullTableName(this.db, this.table)
	condStr := ``
	if condition != nil {
		condStr += ` WHERE ` + condition.GetSql()
	}
	sql := `DELETE FROM ` + ftn + condStr
	return sql
}

func ShowColumns(db, table string) string {
	dt := fullTableName(db, table)
	return `SHOW COLUMNS FROM ` + dt
}

func fullTableName(db, table string) string {
	table = "`" + table + "`"
	if db != `` {
		db = "`" + db + "`"
	}
	return db + "." + table
}
