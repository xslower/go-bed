package orm

import (
	"errors"
	"io/ioutil"

	"github.com/json-iterator/go"
)

var (
	NewConn func(*ConnInfo) IConn = NewMysqlConn
)

const (
	CONFIG_DB_MAPPING = `db-mapping`
)

//db集群连接管理类
//包内使用
type connManager struct {
	conns []IConn
	//分表key与db名的对应
	keyMap map[string]int
	def    IConn
}

func (this *connManager) startWithFile(file string) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	config := ConnConfig{}
	err = jsoniter.Unmarshal(bytes, &config)
	if err != nil {
		return err
	}
	if len(config) == 0 {
		return errors.New(`db config error! no connection info`)
	}
	km := config[CONFIG_DB_MAPPING]
	if km != nil {
		delete(config, CONFIG_DB_MAPPING)
	}
	err = this.start(config, km)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 为db集群建立连接。
 */
func (this *connManager) start(config ConnConfig, km map[string]string) error {
	if len(config) < 1 {
		return errors.New(`db config error! have no connection info`)
	}
	this.keyMap = make(map[string]int)
	this.conns = make([]IConn, 0, len(config))
	//以映射为主，非default也没有映射到的节点，必然用不到，也就无需建立链接
	for key, node := range km {
		ci := config[node]
		if ci == nil {
			return errors.New(`db config error! partition node [` + node + `] not found`)
		}
		conn := NewConn(NewConnInfo(ci))
		if key == `default` {
			this.def = conn
		} else {
			this.conns = append(this.conns, conn)
			this.keyMap[key] = len(this.conns) - 1
		}
	}
	if this.def == nil { //没有配置默认db，则取随机位置作为默认
		if len(this.conns) > 0 {
			this.def = this.conns[0]
		} else {
			for _, ci := range config { //只取第一个所以break
				this.def = NewConn(NewConnInfo(ci))
				if this.def != nil {
					break
				}
			}
		}
	}
	if this.def == nil {
		return errors.New(`db config error! no available db!`)
	}
	return nil
}

func (this *connManager) getConn(key string) (conn IConn) {
	if key == `` || key == `default` {
		return this.def
	}
	idx, ok := this.keyMap[key]
	if !ok {
		return this.def
	}
	conn = this.conns[idx]
	return
}

func NewConnInfo(ci map[string]string) (c *ConnInfo) {
	c = &ConnInfo{}
	c.BindMap(ci)
	return
}

type ConnInfo struct {
	user   string
	pass   string
	host   string
	port   string
	dbname string
}

func (this *ConnInfo) BindMap(ci map[string]string) {
	this.user = ci[`user`]
	this.pass = ci[`pass`]
	this.host = ci[`host`]
	this.port = ci[`port`]
	if this.port == `` {
		this.port = `3306`
	}
	this.dbname = ci[`dbname`]
}

type ConnConfig map[string]map[string]string
