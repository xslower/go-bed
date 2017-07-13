package test

import (
	"github.com/xslower/goutils/utils"
)

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
