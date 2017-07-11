package orm

import (
	"github.com/xslower/go-die/orm/sql"
	"github.com/xslower/goutils/utils"
)

var (
	_conn_manager = &connManager{}
	// fNewSqlBuilder =
)

func getConn(key string) (ic IConn) {
	ic = _conn_manager.getConn(key)
	return
}

func Init(cim map[string]*ConnInfo, km map[string]string) {
	_conn_manager.start(cim, km)
}

func throw(err error, msg ...interface{}) {
	utils.Throw(err, msg...)
}

func echo(i ...interface{}) {
	utils.Echo(i...)
}

func echoStrSlice(strs ...[]string) {
	utils.EchoStrSlice(strs...)
}

func echoBytes(args interface{}) {
	utils.EchoBytes(args)
}
