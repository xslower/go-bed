package builder

var (
	helper = `package {PACKAGE}

import (
	"github.com/xslower/goutils/utils"
)

func throw(err error, msg ...string) {
	utils.Throw(err, msg...)
}

func check(err error, msg ...interface{}) {
	utils.Check(err, msg...)
}

func logit(data ...interface{}) {
	utils.Logit(data...)
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
	`
)
