package orm

import (
	"errors"
)

var (
	NewConn func(*ConnInfo) IConn = NewMysqlConn
)

//db集群连接管理类
//包内使用
type connManager struct {
	conns []IConn
	//分表key与db名的对应
	keyMap map[string]int
	def    IConn
}

func (this *connManager) start(cim map[string]*ConnInfo, km map[string]string) error {
	this.conns = make([]IConn, len(cim))
	var cnt = 0
	//以映射为主，非default也没有映射到的节点，必然用不到，也就无需建立链接
	for key, node := range km {
		ci := cim[node]
		if ci == nil {
			return errors.New(`db config error! partition node [` + node + `] not found`)
		}
		conn := NewConn(ci)
		if key == `default` {
			this.def = conn
		} else {
			this.conns[cnt] = conn
			this.keyMap[key] = cnt
			cnt++
		}
	}
	return nil
}

func (this *connManager) getConn(key string) (conn IConn) {
	idx, ok := this.keyMap[key]
	if !ok {
		return this.def
	}
	conn = this.conns[idx]
	return
}

type ConnInfo struct {
	user   string
	pass   string
	host   string
	port   string
	dbname string
}
