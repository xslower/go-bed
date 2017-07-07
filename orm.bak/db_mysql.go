package orm

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func assembleDsn(dbc map[string]string) string {
	dsn := dbc[`user`] + `:` + dbc[`password`] + `@tcp(` + dbc[`host`] + `:` + dbc[`port`] + `)/` + dbc[`dbname`] + `?charset=utf8`
	return dsn
}

func NewMysqlConn(config map[string]string) IConn {
	mc := &MysqlConn{}
	sqlDb, err := sql.Open(`mysql`, assembleDsn(config))
	throw(err)
	mc.db = sqlDb
	return mc
}

type MysqlConn struct {
	Conn
}

func init() {

}
