package builder

var (
	importDefine = `
import (
	"errors"
	"strconv"
	"strings"
	"github.com/xslower/go-die/orm"
	"hash/crc32"
)
`
	subModelDefine = `
func New{MODEL}() *{MODEL} {
	mt := &{MODEL}{}
	mt.BaseModel.IModel = mt
	mt.result = &([]*{ROW}{})
	return mt
}

type {MODEL} struct {
	orm.BaseModel
	result *[]*{ROW}
}

func (this *{MODEL}) DbName() string {
	return "{DB}"
}
func (this *{MODEL}) TableName(elems ...orm.ISqlElem) string {
	{CODING}
}
func (this *{MODEL}) PartitionKey() []string {
	return []string{"{PARTKEY}"}
}

func (this *{MODEL}) Result() []*{ROW}{
	return *(this.result)
}
func (this *{MODEL}) CreateRow() orm.IRow {
	irow := &{ROW}{}
	*(this.result) = append(*(this.result), irow)
	return irow
}
func (this *{MODEL}) ResetResult() {
	*(this.result) = []*{ROW}{}
}
func (this *{MODEL}) ToIRows(rows *[]*{ROW}) []orm.IRow {
	this.result = rows
	irows := make([]orm.IRow, len(*rows))
	for i, r := range *rows {
		irows[i] = r
	} 
	return irows
}`
	rowSetterDefine = `
func (this *{ROW}) Set(key string, val []byte) error {
	var err error
	switch key {
	{CASE}
	default:
		err = errors.New("No such column [" + key + "]")
	}
	return err
}`
	rowGetterDefine = `
func (this *{ROW}) Get(key string) (interface{}) {
	switch key {
	{CASE}
	default:	
		return nil
	}
}`
	rowColumnsDefine = `
func (this *{ROW}) Columns() []string {
	return []string{{COLUMNS}}
}`

	initDefine = `
func ormInit(file string) (err error) {
	err = orm.InitWithFile(file)
	if err != nil {
		return
	}
	{MODEL}
	return
}

func getPart(idx, part_num int) string {
	pn_str := strconv.Itoa(part_num)
	part := idx % part_num + 1
	// part_str := strconv.FormatUint(uint64(part), 10)
	part_str := strconv.Itoa(part)
	ln := len(pn_str)
	for ln > len(part_str) {
		part_str = "0" + part_str
	}
	return part_str
}

func packageHolder(){
	_ = crc32.ChecksumIEEE([]byte("a"))
	_ = strings.Join([]string{}, "")
}
`

	tbnCodeDefine = `tbn_default := "{TBN_DEFAULT}"
	if len(elems) == 0 {
		return tbn_default
	}
	elem := elems[0]
	tbn := "{TBN_DEFINE}"
	val_default := []string{"{VAL_DEFAULT}"}
	keys := this.PartitionKey()
	for i, k := range keys {
		ifc, _ := elem.Get(k)
		{MODE}
		elem.Del(k)
		tbn = strings.Replace(tbn, "{"+k+"}", val, -1)

	}
	return tbn`
	t1Define = `val, _ := orm.InterfaceToString(ifc)
		if val == "" {
			val = val_default[i]
		}
	`
	t2Define = `idx, _ := orm.InterfaceToInt(ifc)
		val := val_default[i]
		if idx != 0 {
			val = getPart(idx, {PART_NUMBER})
		}
	`
)

// tbl_default := "{DEFAULT}"
// 	if len(elems) == 0 {
// 		return tbl_default
// 	}
// 	elem := elems[0]
// 	ifc, _ := elem.Get("{FIELD}")

// val, _ := orm.InterfaceToString(ifc)
// 	elem.Del("{FIELD}")
// 	if val == "" {
// 		return "{PREFIX}{DEFAULT}"
// 	}
// 	return "{PREFIX}" + val

// val, _ := orm.InterfaceToInt(ifc)
// 	if val == 0 {
// 		return "{PREFIX}{DEFAULT}"
// 	}
// 	part := val % {PART_NUMBER}
// 	prefix := "{PREFIX}"
// 	if part < 10 {
// 		prefix += "0"
// 	}
// 	return prefix + strconv.FormatUint(uint64(part), 10)

// val, _ := orm.InterfaceToString(ifc)
// 	if val == "" {
// 		return "{PREFIX}{DEFAULT}"
// 	}
// 	hash := {HASH}
// 	part := hash % {PART_NUMBER}
// 	prefix := "{PREFIX}"
// 	if part < 10 {
// 		prefix += "0"
// 	}
// 	return prefix + strconv.FormatUint(uint64(part), 10)
//
