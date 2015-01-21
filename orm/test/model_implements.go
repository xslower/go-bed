package test

import (
	"errors"
	"greentea/orm"
	"hash/crc32"
	"strconv"
)

func NewMytestModel() *MytestModel {
	mt := &MytestModel{}
	mt.BaseModel.IModel = mt
	mt.result = []*Mytest{}
	return mt
}

type MytestModel struct {
	orm.BaseModel
	result []*Mytest
}

func (this *MytestModel) DbName() string {
	return "test"
}
func (this *MytestModel) TableName(elem orm.ISqlElem) string {
	ifc := elem.Get("category")
	category := orm.InterfaceToString(ifc)
	if category == "" {
		return "mytest_tvplay"
	}
	return "mytest_" + category
}

func (this *MytestModel) Result() []*Mytest {
	return this.result
}
func (this *MytestModel) CreateRow() orm.IRow {
	irow := &Mytest{}
	this.result = append(this.result, irow)
	return irow
}
func (this *MytestModel) ResetResult() {
	this.result = []*Mytest{}
}

func (this *Mytest) Set(key string, val []byte) error {
	var err error
	switch key {
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

func NewPhpTestModel() *PhpTestModel {
	mt := &PhpTestModel{}
	mt.BaseModel.IModel = mt
	mt.result = []*PhpTest{}
	return mt
}

type PhpTestModel struct {
	orm.BaseModel
	result []*PhpTest
}

func (this *PhpTestModel) DbName() string {
	return "test"
}
func (this *PhpTestModel) TableName(elem orm.ISqlElem) string {
	ifc := elem.Get("id")
	id := orm.InterfaceToInt(ifc)
	if id == 0 {
		return "php_test_01"
	}
	part := id % 100
	prefix := "php_test_"
	if part < 10 {
		prefix += "0"
	}
	return prefix + strconv.FormatUint(uint64(part), 10)
}

func (this *PhpTestModel) Result() []*PhpTest {
	return this.result
}
func (this *PhpTestModel) CreateRow() orm.IRow {
	irow := &PhpTest{}
	this.result = append(this.result, irow)
	return irow
}
func (this *PhpTestModel) ResetResult() {
	this.result = []*PhpTest{}
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

func NewUserModel() *UserModel {
	mt := &UserModel{}
	mt.BaseModel.IModel = mt
	mt.result = []*User{}
	return mt
}

type UserModel struct {
	orm.BaseModel
	result []*User
}

func (this *UserModel) DbName() string {
	return ""
}
func (this *UserModel) TableName(elem orm.ISqlElem) string {
	ifc := elem.Get("name")
	name := orm.InterfaceToString(ifc)
	if name == "" {
		return "php_test_01"
	}
	hash := crc32.ChecksumIEEE([]byte(name))
	part := hash % 100
	prefix := "php_test_"
	if part < 10 {
		prefix += "0"
	}
	return prefix + strconv.FormatUint(uint64(part), 10)
}

func (this *UserModel) Result() []*User {
	return this.result
}
func (this *UserModel) CreateRow() orm.IRow {
	irow := &User{}
	this.result = append(this.result, irow)
	return irow
}
func (this *UserModel) ResetResult() {
	this.result = []*User{}
}

func (this *User) Set(key string, val []byte) error {
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

func (this *User) Get(key string) interface{} {
	switch key {
	case "id":
		return this.Id
	case "name":
		return this.Name

	default:
		return nil
	}
}

func ormStart(dbConfig map[string]string) {
	orm.Start(dbConfig)
	orm.RegisterModel(NewMytestModel())
	orm.RegisterModel(NewPhpTestModel())
	orm.RegisterModel(NewUserModel())

}
func packageHolder() {
	_ = crc32.ChecksumIEEE([]byte("a"))
}
