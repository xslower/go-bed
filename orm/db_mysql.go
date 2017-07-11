package orm

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// func assembleDsn(dbc map[string]string) string {
// 	dsn := dbc[`username`] + `:` + dbc[`password`] + `@tcp(` + dbc[`host`] + `:` + dbc[`port`] + `)/` + dbc[`dbname`] + `?charset=utf8`
// 	return dsn
// }
func assembleDsn(ci *ConnInfo) (dsn string) {
	dsn = ci.user + `:` + ci.pass + `@tcp(` + ci.host + `:` + ci.port + `)/` + ci.dbname + `?charset=utf8`
	return
}
func NewMysqlConn(ci *ConnInfo) IConn {
	mc := &MysqlConn{}
	sqlDb, err := sql.Open(`mysql`, assembleDsn(ci))
	throw(err)
	mc.db = sqlDb
	return mc
}

type MysqlConn struct {
	Conn
}

func init() {

}
