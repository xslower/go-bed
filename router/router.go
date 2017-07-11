package router

import (
	// `encoding/json`
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/xslower/goutils/json"
)

var (
	err_not_found = errors.New(`not found`)
)

func init() {
	// _handler_map[`ctrl`] = control
	// http.HandleFunc(`/`, httpEntrance)

}

//
type Context struct {
	w http.ResponseWriter
	r *http.Request
	// 临时数据
	td interface{}
}

func (c *Context) StoreData(d interface{}) {
	c.td = d
}

func (c *Context) FetchData() interface{} {
	return c.td
}

func (c *Context) Print(s ...interface{}) {
	fmt.Fprint(c.w, s...)
}

func (c *Context) SendJson(data interface{}) {
	// var bytes, err = json.Marshal(data)
	var jn = json.Encode(data)
	// throw(err, `result to json`)
	c.Print(`{"err_no":0,"result":`, jn, `}`)
}

func (c *Context) SendError(msg interface{}) {
	c.Print(`{"err_no":1,"err_msg":"`, msg, `"}`)
}

type BeforeHandler func(*Context) error
type HttpHandler func(*Context) (interface{}, error)
type Event string

// type AfterHandler func(*Context)

func New() *Router {
	return &Router{
		beforeMap:  make(map[string][]BeforeHandler),
		handlerMap: make(map[string]HttpHandler),
	}
}

type Router struct {
	//例如权限验证、ip限制、参数过滤等预处理
	beforeMap  map[string][]BeforeHandler
	handlerMap map[string]HttpHandler
}

func (this *Router) BindBeforeHandler(path string, bh BeforeHandler) {
	var bhs = this.beforeMap[path]
	if bhs == nil {
		bhs = make([]BeforeHandler, 5)
		bhs[0] = bh
		this.beforeMap[path] = bhs
	} else {
		bhs = append(bhs, bh)
		this.beforeMap[path] = bhs
	}

}

func (this *Router) BindHandler(path string, hh HttpHandler) {
	this.handlerMap[path] = hh
}

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ctx = &Context{w: w, r: r}
	defer func() {
		if e := recover(); e != nil {
			ctx.SendError(e)
		}
	}()
	var uri = r.URL.RequestURI()[1:]
	var pos = strings.Index(uri, `?`)
	if pos < 0 {
		pos = len(uri)
	}
	var path = uri[:pos]
	var handler, ok = this.handlerMap[path]
	if !ok {
		ctx.SendError(err_not_found)
		return
	}
	r.ParseForm()
	//进行预处理
	bhs, ok := this.beforeMap[path]
	stop := false
	if ok {
		for _, bh := range bhs {
			err := bh(ctx)
			if err != nil {
				ctx.SendError(err)
				stop = true
				break
			}
		}
	}
	if stop {
		return
	}
	var ret, err = handler(ctx)
	if err != nil {
		ctx.SendError(err)
		return
	}
	ctx.SendJson(ret)

}

func control(c *Context) interface{} {
	return `control`
}
