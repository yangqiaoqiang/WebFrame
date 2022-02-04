package gin

//框架初始逻辑构造 基于package http
import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

//GET method 增加get 请求
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

//POST method 增加post请求
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

//Run method 启动httpserver
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

//重写handle接口方法
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND:%s\n", req.URL)
	}
}
