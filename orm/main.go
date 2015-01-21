package orm

import (
	"fmt"
	`log`
)

func init() {

}

func Start(dbCnf map[string]string) {
	gDbRegistry[`mysql`] = NewMysqlConn
	gSqlBuFnRegistry[`mysql`] = NewMysqlBuilder
	driver, ok := dbCnf[`driver`]
	if !ok { //default driver is mysql
		driver = `mysql`
	}
	dbFunc, ok := gDbRegistry[driver]
	if !ok {
		panic(`Do not support driver: [` + driver + `] in db driver`)
	}
	gDbConn = dbFunc(dbCnf)

	fNewSqlBuilder, ok = gSqlBuFnRegistry[driver]
	if !ok {
		panic(`Do not support driver: [` + driver + `] in sql builder`)
	}
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
