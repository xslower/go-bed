package orm

import ()

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
