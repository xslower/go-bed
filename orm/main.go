package orm

import (
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

//使用配置文件初始化orm
func InitWithFile(file string) (err error) {
	err = _conn_manager.startWithFile(file)
	return
}

func Init(cim ConnConfig, km map[string]string) (err error) {
	err = _conn_manager.start(cim, km)
	return
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
