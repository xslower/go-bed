package test

import (
	"errors"
	"github.com/xslower/go-die/orm"
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
	return "mytest"
}
func (this *MytestModel) PartitionKey() []string {
	return []string{""}
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

func ormInit(file string) (err error) {
	err = orm.InitWithFile(file)
	if err != nil {
		return
	}
	orm.RegisterModel(NewMytestModel())

	return
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
	_ = strings.Join([]string{}, "")
}
