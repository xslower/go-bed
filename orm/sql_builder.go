/**
* sql builder for mysql
 */
package orm

import (
	"strings"

	"github.com/xslower/goutils/utils"
)

var (
	NewSqlBuilder = NewMysqlBuilder
)

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
		fldStr = strings.Trim(fields.GetSql(), ` `)
		if fldStr == `` {
			fldStr = `*`
		}
	}
	condStr := this.getCondStr(condition)

	appendStr := ``
	if appends != nil {
		appendStr = appends.GetSql()
	}
	sql := `SELECT ` + fldStr + ` FROM ` + ftn + condStr + ` ` + appendStr
	return sql
}

func (this *MysqlBuilder) getFieldWriteStr(fields IFields) string {
	if fields == nil {
		panic(`the [fields] is nil`)
	}
	str := strings.Trim(fields.GetSql(), ` `)
	if str == `` {
		panic(`the [fields] is empty`)
	}
	return str
}

func (this *MysqlBuilder) getCondStr(cond ICondition) string {
	if cond == nil {
		return ``
	}
	str := strings.Trim(cond.GetSql(), ` `)
	if str != `` {
		str = ` WHERE ` + str
	}
	return str
}

func (this *MysqlBuilder) MultiInsert(fields IFields) string {
	fldStr := this.getFieldWriteStr(fields)
	ftn := fullTableName(this.db, this.table)
	sql := `INSERT ` + ftn + ` ` + fldStr
	return sql
}

func (this *MysqlBuilder) Insert(fields IFields) string {
	fldStr := this.getFieldWriteStr(fields)
	ftn := fullTableName(this.db, this.table)
	sql := `INSERT ` + ftn + ` SET ` + fldStr
	return sql
}

func (this *MysqlBuilder) InsertOrUpdate(fields IFields) string {
	fldStr := this.getFieldWriteStr(fields)
	ftn := fullTableName(this.db, this.table)
	sql := `INSERT ` + ftn + ` SET ` + fldStr + ` ON DUPLICATE KEY UPDATE ` + fldStr
	return sql
}

func (this *MysqlBuilder) Update(fields IFields, condition ICondition) string {
	fldStr := this.getFieldWriteStr(fields)
	ftn := fullTableName(this.db, this.table)
	condStr := this.getCondStr(condition)
	sql := `UPDATE ` + ftn + ` SET ` + fldStr + condStr
	return sql
}

func (this *MysqlBuilder) Delete(condition ICondition) string {
	ftn := fullTableName(this.db, this.table)
	condStr := this.getCondStr(condition)
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
		db = "`" + db + "`."
	}
	return db + table
}

func IfcToSqlValue(ifc interface{}) string {
	val, typ := utils.InterfaceToString(ifc)
	if typ == utils.TYPE_STRING {
		return `'` + val + `'`
	}
	return val
}
