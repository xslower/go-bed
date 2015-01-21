package test

import (
	`fmt`
	`github.com/astaxie/beego/config`
	`greentea/orm`
	`log`
	"testing"
)

var (
	gid = []int{}
	gmt *MytestModel
)

func init() {
	cnf, err := config.NewConfig(`ini`, `config.ini`)
	throw(err)
	dbCnf, err := cnf.GetSection(`database.147`)
	throw(err)
	ormStart(dbCnf)
	gmt = NewMytestModel()
}

func TestInsert(t *testing.T) {
	fm := map[string]interface{}{`name`: `这是单条插入测试`}
	f := orm.NewFieldsWriteWithMap(fm)
	id := gmt.Insert(f)
	gid = append(gid, id)
}

func TestMultiInsert(t *testing.T) {
	fm := map[string]interface{}{`name`: `这是多条插入测试1`}
	f := orm.NewFieldsMultiInsertWithMap(fm)
	f.AddRow(`这是多条插入测试2`)
	gmt.MultiInsert(f)
}
func TestUpdate(t *testing.T) {
	fm := map[string]interface{}{`name`: `修改后的名称`}
	f := orm.NewFieldsWriteWithMap(fm)
	cond := orm.NewConds()
	cond.Equal(`id`, gid[0])
	gmt.Update(f, cond)
}
func TestSelect(t *testing.T) {

}

func final() {
	if exception := recover(); exception != nil {
		log.Println(exception)
	}
}
func throw(err error) {
	if err != nil {
		panic(err.Error())
	}
}
func check(err error) bool {
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
func echo(i ...interface{}) {
	fmt.Println(i...)
}
func echoStrings(strs ...[]string) {
	for i, val := range strs {
		fmt.Print(`[`, i, `] `)
		for j, v := range val {
			fmt.Println(j, `: `, v)
		}
		fmt.Print("\n")
	}
}
func echoBytes(args interface{}) {
	switch v := args.(type) {
	case []byte:
		echo(string(v))
	case []rune:
		echo(string(v))
	case [][]byte:
		for i, val := range v {
			fmt.Print(i, `: `, string(val), ` `)
		}
		fmt.Print("\n")
	case [][]rune:
		for i, val := range v {
			fmt.Print(i, `: `, string(val), ` `)
		}
		fmt.Print("\n")
	case [][][]byte:
		for i, val := range v {
			fmt.Print(i, ` `)
			echoBytes(val)
		}
	case [][][]rune:
		for i, val := range v {
			fmt.Print(i, ` `)
			echoBytes(val)
		}
	default:
		echo(v)
	}

}

func logit(data ...interface{}) {
	log.Println(data...)
}
