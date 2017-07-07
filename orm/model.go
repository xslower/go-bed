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
	PartitionKey() []string //分表的字段，如果没有则返回空
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

func (this *BaseModel) PartitionKey() []string {
	return this.IModel.PartitionKey()
}

func (this *BaseModel) Insert(fields IFields) int {
	if fields == nil {
		panic(`fields is nil, expect IFields`)
	}
	sb := fNewSqlBuilder(this.DbName(), this.TableName(fields))
	this.filterColumn(fields)
	sql := sb.Insert(fields)
	values := fields.Values()
	this.lastSql = sql
	id := gDbConn.Exec(sql, values...)
	return id
}

func (this *BaseModel) InsertOrUpdate(fields IFields) int {
	if fields == nil {
		panic(`fields is nil, expect IFields`)
	}
	sb := fNewSqlBuilder(this.DbName(), this.TableName(fields))
	this.filterColumn(fields)
	sql := sb.InsertOrUpdate(fields)
	values := fields.Values()
	values = append(values, values...) //这里insert和update需要两套值
	this.lastSql = sql
	id := gDbConn.Exec(sql, values...)
	return id
}

func (this *BaseModel) MultiInsert(fields IFields) int {
	if fields == nil {
		panic(`fields is nil, expect IFields`)
	}
	sb := fNewSqlBuilder(this.DbName(), this.TableName(fields))
	// this.filterColumn(fields)
	sql := sb.MultiInsert(fields)
	values := fields.Values()
	this.lastSql = sql
	id := gDbConn.Exec(sql, values...)
	return id
}

func (this *BaseModel) Delete(condition ICondition) int {
	if condition == nil {
		panic(`condition is nil, expect ICondition`)
	}
	sb := fNewSqlBuilder(this.DbName(), this.TableName(condition))
	this.filterColumn(condition)
	if len(condition.Columns()) == 0 {
		panic(`do not allow update table without any condition!`)
	}
	sql := sb.Delete(condition)
	values := condition.Values()
	this.lastSql = sql
	num := gDbConn.Exec(sql, values...)
	return num
}

func (this *BaseModel) Update(fields IFields, condition ICondition) int {
	if fields == nil {
		panic(`fields is nil, expect IFields`)
	}
	if condition == nil { //patch update table is a dangerous action
		panic(`condition is nil, expect ICondition`)
	}
	sb := fNewSqlBuilder(this.DbName(), this.TableName(condition))
	this.filterColumn(fields)
	this.filterColumn(condition)
	if len(condition.Columns()) == 0 {
		panic(`do not allow update table without any condition!`)
	}
	sql := sb.Update(fields, condition)
	// echo(fields.Values())
	// echo(condition.Values())
	values := append(fields.Values(), condition.Values()...)
	this.lastSql = sql
	num := gDbConn.Exec(sql, values...)
	return num
}

func (this *BaseModel) Select(condition ICondition, appends IAppends, fields IFields) bool {
	if condition == nil {
		condition = NewConds()
	}
	if appends == nil {
		appends = NewAppends()
	}
	if fields == nil {
		fields = NewFieldsRead()
	}
	sb := fNewSqlBuilder(this.DbName(), this.TableName(condition))
	this.filterColumn(fields)
	this.filterColumn(condition)
	sql := sb.Select(fields, condition, appends)
	// values := append(fields.Values(), condition.Values()...)
	values := append(condition.Values(), appends.Values()...)
	this.lastSql = sql
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

func (this *BaseModel) Read(irow IRow) bool {
	if irow == nil {
		return false
	}
	iConds := this.IRowToConds(irow)
	sb := fNewSqlBuilder(this.DbName(), this.TableName(iConds))
	this.filterColumn(iConds) //有些model中的分表字段可能不在表中
	iApds := NewAppends().Limit(0, 1)
	values := append(iConds.Values(), iApds.Values()...)
	sql := sb.Select(nil, iConds, iApds)
	this.lastSql = sql
	columns, rrows := gDbConn.Query(sql, values...)
	if len(rrows) > 0 {
		for i, clm := range columns {
			irow.Set(clm, rrows[0][i])
		}
		return true
	}
	return false
}

func (this *BaseModel) ReadMulti(result []IRow, cond IRow) bool {
	iConds := this.IRowToConds(cond)
	return this.Select(iConds, nil, nil)
}

func (this *BaseModel) Save(irow IRow) int {
	iFields := this.IRowToFields(irow)
	id := this.InsertOrUpdate(iFields)
	return id
}

func (this *BaseModel) SaveMulti(irows []IRow) int {
	if len(irows) == 0 {
		return 0
	}
	fields := this.IRowsToFieldsMI(irows)
	if len(fields.Columns()) == 0 {
		return 0
	}
	return this.MultiInsert(fields)
}

func (this *BaseModel) Remove(irow IRow) int {
	ti := this.getTableInfo()
	pks := ti.primaryKey
	partKey := this.PartitionKey()
	pks = append(pks, partKey...)
	condCol, condVal := this.getRowValues(irow, pks)
	iConds := NewConds().InitWithSlice(condCol, condVal)
	num := this.Delete(iConds)
	return num
}

func (this *BaseModel) IRowToConds(irow IRow) *Condition {
	if irow == nil {
		return NewConds()
	}
	columns := this.getAllColumns()
	columns, values := this.getRowValues(irow, columns)
	conds := NewConds().InitWithSlice(columns, values)
	return conds
}

func (this *BaseModel) IRowToFields(irow IRow) *FieldsWrite {
	if irow == nil {
		return NewFieldsWrite()
	}
	columns := this.getAllColumns()
	columns, values := this.getRowValues(irow, columns)
	fields := NewFieldsWrite().InitWithSlice(columns, values)
	return fields
}

func (this *BaseModel) IRowsToFieldsMI(irows []IRow) *FieldsMultiInsert {
	columns := this.getAllColumns()
	fields := NewFieldsMultiInsert()
	for _, ir := range irows {
		if ir == nil {
			continue
		}
		cols, vals := this.getRowValues(ir, columns)
		fields.AddSlice(vals, cols...)
	}
	return fields
}

/**
 * 指定字段名slice获取字段值slice
 */
func (this *BaseModel) getRowValues(irow IRow, columns []string) (newColumns []string, values []interface{}) {
	// negligible := map[string]bool{}
	// if trimNegligible {
	// 	ti := this.getTableInfo()
	// 	negColumns := ti.Negligible()
	// 	for _, clm := range negColumns {
	// 		negligible[clm] = true
	// 	}
	// }
	// newColumns := []string{}
	for _, clm := range columns {
		ifc := irow.Get(clm)
		val, _ := InterfaceToString(ifc)
		if val == `0` || val == `` { //去除0值，因为OO式的db操作0值没有意义
			continue
		}
		// if trimNegligible && negligible[clm] { //去除可忽略
		// 	continue
		// }
		newColumns = append(newColumns, clm)
		values = append(values, ifc)
	}
	return
}

func (this *BaseModel) getAllColumns() []string {
	ti := this.getTableInfo()
	columns := ti.Columns()
	partKey := this.PartitionKey()
	return append(columns, partKey...)
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
