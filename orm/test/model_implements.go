package test

import (
	"errors"
	"greentea/orm"
	"hash/crc32"
	"strconv"
	"strings"
)

func NewMytestModel() *MytestModel {
	mt := &MytestModel{}
	mt.BaseModel.IModel = mt
	mt.result = &([]*Mytest{})
	return mt
}

type MytestModel struct {
	orm.BaseModel
	result *[]*Mytest
}

func (this *MytestModel) DbName() string {
	return "test"
}
func (this *MytestModel) TableName(elems ...orm.ISqlElem) string {
	tbn_default := "mytest_d"
	if len(elems) == 0 {
		return tbn_default
	}
	elem := elems[0]
	tbn := "my{test}_{category}"
	val_default := []string{"test", "d"}
	keys := this.PartitionKey()
	for i, k := range keys {
		ifc, _ := elem.Get(k)
		val, _ := orm.InterfaceToString(ifc)
		if val == "" {
			val = val_default[i]
		}

		elem.Del(k)
		tbn = strings.Replace(tbn, "{"+k+"}", val, -1)

	}
	return tbn
}
func (this *MytestModel) PartitionKey() []string {
	return []string{"test", "category"}
}

func (this *MytestModel) Result() []*Mytest {
	return *(this.result)
}
func (this *MytestModel) CreateRow() orm.IRow {
	irow := &Mytest{}
	*(this.result) = append(*(this.result), irow)
	return irow
}
func (this *MytestModel) ResetResult() {
	*(this.result) = []*Mytest{}
}
func (this *MytestModel) ToIRows(rows *[]*Mytest) []orm.IRow {
	this.result = rows
	irows := make([]orm.IRow, len(*rows))
	for i, r := range *rows {
		irows[i] = r
	}
	return irows
}

func (this *Mytest) Set(key string, val []byte) error {
	var err error
	switch key {
	case "test":
		this.Test = string(val)
	case "category":
		this.Category = string(val)
	case "id":
		this.Id, err = strconv.Atoi(string(val))
	case "name":
		this.Name = string(val)
	case "create_time":
		this.CreateTime = string(val)

	default:
		err = errors.New("No such column [" + key + "]")
	}
	return err
}

func (this *Mytest) Get(key string) interface{} {
	switch key {
	case "test":
		return this.Test
	case "category":
		return this.Category
	case "id":
		return this.Id
	case "name":
		return this.Name
	case "create_time":
		return this.CreateTime

	default:
		return nil
	}
}

func (this *Mytest) Columns() []string {
	return []string{`test`, `category`, `id`, `name`, `create_time`}
}

func NewPhpTestModel() *PhpTestModel {
	mt := &PhpTestModel{}
	mt.BaseModel.IModel = mt
	mt.result = &([]*PhpTest{})
	return mt
}

type PhpTestModel struct {
	orm.BaseModel
	result *[]*PhpTest
}

func (this *PhpTestModel) DbName() string {
	return "test"
}
func (this *PhpTestModel) TableName(elems ...orm.ISqlElem) string {
	tbn_default := "php_test_001"
	if len(elems) == 0 {
		return tbn_default
	}
	elem := elems[0]
	tbn := "php_test_{id}"
	val_default := []string{"001"}
	keys := this.PartitionKey()
	for i, k := range keys {
		ifc, _ := elem.Get(k)
		idx, _ := orm.InterfaceToInt(ifc)
		val := val_default[i]
		if idx != 0 {
			val = getPart(idx, 100)
		}

		elem.Del(k)
		tbn = strings.Replace(tbn, "{"+k+"}", val, -1)

	}
	return tbn
}
func (this *PhpTestModel) PartitionKey() []string {
	return []string{"id"}
}

func (this *PhpTestModel) Result() []*PhpTest {
	return *(this.result)
}
func (this *PhpTestModel) CreateRow() orm.IRow {
	irow := &PhpTest{}
	*(this.result) = append(*(this.result), irow)
	return irow
}
func (this *PhpTestModel) ResetResult() {
	*(this.result) = []*PhpTest{}
}
func (this *PhpTestModel) ToIRows(rows *[]*PhpTest) []orm.IRow {
	this.result = rows
	irows := make([]orm.IRow, len(*rows))
	for i, r := range *rows {
		irows[i] = r
	}
	return irows
}

func (this *PhpTest) Set(key string, val []byte) error {
	var err error
	switch key {
	case "id":
		this.Id, err = strconv.Atoi(string(val))
	case "name":
		this.Name = string(val)

	default:
		err = errors.New("No such column [" + key + "]")
	}
	return err
}

func (this *PhpTest) Get(key string) interface{} {
	switch key {
	case "id":
		return this.Id
	case "name":
		return this.Name

	default:
		return nil
	}
}

func (this *PhpTest) Columns() []string {
	return []string{`id`, `name`}
}

func ormStart(dbConfig map[string]string) {
	orm.Start(dbConfig)
	orm.RegisterModel(NewMytestModel())
	orm.RegisterModel(NewPhpTestModel())

}

func getPart(idx, part_num int) string {
	pn_str := strconv.Itoa(part_num)
	part := idx%part_num + 1
	// part_str := strconv.FormatUint(uint64(part), 10)
	part_str := strconv.Itoa(part)
	ln := len(pn_str)
	for ln > len(part_str) {
		part_str = "0" + part_str
	}
	return part_str
}

func packageHolder() {
	_ = crc32.ChecksumIEEE([]byte("a"))
}
