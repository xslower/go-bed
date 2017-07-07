package orm

import (
	// `database/sql`
	`reflect`
	`strconv`
	`strings`
)

var (
	//BaseModel不能跟TableInfo绑定，因为tableInfo确定只能有一个，而Model因为为了简化多线程编程，可以有多个。不然result会被争用，导致数据混乱
	gModelRegistry = map[string]*TableInfo{}
)

type TableInfo struct {
	dbName     string
	tableName  string
	primaryKey []string
	columns    []*ColumnInfo
}

func (this *TableInfo) Columns() []string {
	columns := []string{}
	for _, ci := range this.columns {
		columns = append(columns, ci.name)
	}
	return columns
}

//插入时可以忽略的键
func (this *TableInfo) Negligible() []string {
	columns := []string{}
	for _, ci := range this.columns {
		if ci.autoInc || ci.timestamp {
			columns = append(columns, ci.name)
		}
	}
	return columns
}

type ColumnInfo struct {
	name      string
	clmType   uint8
	clmLen    uint16
	autoInc   bool
	timestamp bool //mean whether could skiped this column on insert
}

func RegisterModel(model interface{}) {
	val := reflect.ValueOf(model)

	tnFunc := val.MethodByName(`TableName`)
	tnv := tnFunc.Call([]reflect.Value{})
	tableName := tnv[0].String()
	if tableName == `` {
		panic(`Can't fetch table name!`)
	}
	dnFunc := val.MethodByName(`DbName`)
	dnv := dnFunc.Call([]reflect.Value{})
	dbName := dnv[0].String()
	if gDbConn == nil {
		panic(`You must trigger orm.Start() before orm.RegisterModel()!`)
	}
	_, rrows := gDbConn.Query(ShowColumns(dbName, tableName))
	tblInfo := &TableInfo{dbName: dbName, tableName: tableName}
	for _, rr := range rrows {
		colInfo := &ColumnInfo{}
		//parse Field
		colInfo.name = string(rr[0])
		//parse Type
		tmpType := string(rr[1])
		pos := strings.Index(tmpType, `(`)
		upos := strings.Index(tmpType, `unsigned`)
		if pos > 0 {
			colType := tmpType[:pos]
			if upos > 0 {
				colType = `u` + colType
			}
			upos = strings.Index(tmpType, `)`)
			colLen, err := strconv.Atoi(tmpType[pos+1 : upos])
			if err == nil {
				colInfo.clmLen = uint16(colLen)
			}
			colInfo.clmType = gDbTypeR[colType]
		} else {
			colInfo.clmType = gDbTypeR[tmpType]
		}
		//rr[2] is Null skip
		//parse Key
		if string(rr[3]) == `PRI` {
			tblInfo.primaryKey = append(tblInfo.primaryKey, colInfo.name)
		}
		//parse Default
		tmpDefault := string(rr[4])
		if tmpDefault == `CURRENT_TIMESTAMP` {
			colInfo.timestamp = true
		}
		//parse Extra
		tmpExtra := string(rr[5])
		if tmpExtra == `auto_increment` {
			colInfo.autoInc = true
		}
		//parse over

		tblInfo.columns = append(tblInfo.columns, colInfo)
	}
	gModelRegistry[dbName+tableName] = tblInfo
}

type IModel interface {
	DbName() string
	TableName(elems ...ISqlElem) string
	PartitionKey() string //分表的字段，如果没有则返回空
	CreateRow() IRow
	ResetResult()
}

type IRow interface {
	Set(key string, val []byte) error
	Get(key string) interface{}
	Columns() []string
}

type BaseModel struct {
	lastSql string
	IModel
}

func (this *BaseModel) DbName() string {
	return this.IModel.DbName()
}

func (this *BaseModel) TableName(elems ...ISqlElem) string {
	return this.IModel.TableName(elems...)
}

func (this *BaseModel) PartitionKey() string {
	return this.IModel.PartitionKey()
}

func (this *BaseModel) Insert(fields IFields) int {
	sb := fNewSqlBuilder(this.DbName(), this.TableName(fields))
	this.filterColumn(fields)
	sql := sb.Insert(fields)
	values := fields.Values()
	id := gDbConn.Exec(sql, values...)
	this.lastSql = sql
	return id
}

func (this *BaseModel) InsertOrUpdate(fields IFields) int {
	sb := fNewSqlBuilder(this.DbName(), this.TableName(fields))
	this.filterColumn(fields)
	sql := sb.InsertOrUpdate(fields)
	values := fields.Values()
	values = append(values, values...) //这里insert和update需要两套值
	id := gDbConn.Exec(sql, values...)
	this.lastSql = sql
	return id
}

func (this *BaseModel) MultiInsert(fields IFields) {
	sb := fNewSqlBuilder(this.DbName(), this.TableName(fields))
	this.filterColumn(fields)
	sql := sb.MultiInsert(fields)
	values := fields.Values()
	gDbConn.Exec(sql, values...)
	this.lastSql = sql
}

func (this *BaseModel) Delete(condition ICondition) int {
	sb := fNewSqlBuilder(this.DbName(), this.TableName(condition))
	this.filterColumn(condition)
	sql := sb.Delete(condition)
	values := condition.Values()
	num := gDbConn.Exec(sql, values...)
	this.lastSql = sql
	return num
}

func (this *BaseModel) Update(fields IFields, condition ICondition) int {
	sb := fNewSqlBuilder(this.DbName(), this.TableName(condition))
	this.filterColumn(fields)
	this.filterColumn(condition)
	sql := sb.Update(fields, condition)
	// echo(fields.Values())
	// echo(condition.Values())
	values := append(fields.Values(), condition.Values()...)
	num := gDbConn.Exec(sql, values...)
	this.lastSql = sql
	return num
}

func (this *BaseModel) Select(fields IFields, condition ICondition, appends IAppends) bool {
	sb := fNewSqlBuilder(this.DbName(), this.TableName(condition))
	this.filterColumn(fields)
	this.filterColumn(condition)
	sql := sb.Select(fields, condition, appends)
	values := append(fields.Values(), condition.Values()...)
	values = append(values, appends.Values()...)
	columns, rrows := gDbConn.Query(sql, values...)
	this.IModel.ResetResult()
	// irowSlice := []IRow{}
	for _, rr := range rrows {
		irow := this.IModel.CreateRow()
		for i, clm := range columns {
			irow.Set(clm, rr[i])
		}
		// irowSlice = append(irowSlice, irow)
	}
	return true
}

func (this *BaseModel) ObjRead(irow IRow) bool {
	ti := this.getTableInfo()
	pks := ti.primaryKey
	// condMap := map[string]interface{}{}
	condVal := this.getRowValues(&pks, irow, false)
	iConds := NewCondsWithSlice(pks, condVal)
	sb := fNewSqlBuilder(this.DbName(), this.TableName(iConds))
	this.filterColumn(iConds)
	iApds := NewAppends().Limit(0, 1)
	sql := sb.Select(nil, iConds, iApds)
	columns, rrows := gDbConn.Query(sql, iConds.Values()...)
	if len(rrows) > 0 {
		for i, clm := range columns {
			irow.Set(clm, rrows[0][i])
		}
		return true
	}
	return false
}

func (this *BaseModel) ObjSave(irow IRow) int {
	ti := this.getTableInfo()
	columns := ti.Columns()
	// columns := irow.Columns()
	values := this.getRowValues(&columns, irow, true)

	iFields := NewFieldsWriteWithSlice(columns, values)
	id := this.InsertOrUpdate(iFields)
	return id
}

func (this *BaseModel) ObjRemove(irow IRow) int {
	ti := this.getTableInfo()
	pks := ti.primaryKey
	condVal := this.getRowValues(&pks, irow, false)
	iConds := NewCondsWithSlice(pks, condVal)
	num := this.Delete(iConds)
	return num
}

func (this *BaseModel) getRowValues(columns *[]string, irow IRow, trimNegligible bool) []interface{} {
	negligible := map[string]bool{}
	if trimNegligible {
		ti := this.getTableInfo()
		negColumns := ti.Negligible()
		for _, clm := range negColumns {
			negligible[clm] = true
		}
	}
	clms := []string{}
	partKey := this.PartitionKey()
	hasPartKey := false
	// condMap := map[string]interface{}{}
	values := []interface{}{}
	for _, clm := range *columns {
		if clm == partKey {
			hasPartKey = true
		}
		ifc := irow.Get(clm)
		if trimNegligible && negligible[clm] {
			val, _ := InterfaceToString(ifc)
			if val == `0` || val == `` {
				continue
			}
		}
		clms = append(clms, clm)
		values = append(values, ifc)
	}
	if !hasPartKey && partKey != `` {
		clms = append(clms, partKey)
		values = append(values, irow.Get(partKey))
	}
	*columns = clms
	return values
}

func (this *BaseModel) LastSql() string {
	return this.lastSql
}

func (this *BaseModel) filterColumn(elem ISqlElem) {
	ti := this.getTableInfo()
	exist := map[string]bool{}
	for _, clm := range ti.Columns() {
		exist[clm] = true
	}
	columns := elem.Columns()
	for _, clm := range columns {
		if !exist[clm] { //此字段不在本表中
			elem.Del(clm)
		}
		if ifc, _ := elem.Get(clm); ifc == nil { //如果值为nil则去除
			elem.Del(clm)
		}
	}
}

func (this *BaseModel) getTableInfo() *TableInfo {
	tiKey := this.DbName() + this.TableName()
	ti, ok := gModelRegistry[tiKey]
	if !ok {
		panic(`Not registed such Model:[` + tiKey + `]`)
	}
	return ti
}
