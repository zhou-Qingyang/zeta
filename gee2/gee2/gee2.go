package gee2

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Engine struct {
	router *Router
	*RouterGroup
	groups []*RouterGroup
}

func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 查找需要使用的handlers
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			fmt.Println("find middleware", group.prefix)
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := NewContext(w, r)
	c.handlers = middlewares
	e.router.handle(c)
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) GET(path string, handler HandlerFunc) {
	e.router.addRoute(http.MethodGet, path, handler)
}

func (e *Engine) POST(path string, handler HandlerFunc) {
	e.router.addRoute(http.MethodPost, path, handler)
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	group := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: engine,
	}
	g.engine.groups = append(g.engine.groups, group)
	return group
}

func (g *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := g.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	g.engine.router.addRoute(method, pattern, handler)
}

func (g *RouterGroup) GET(path string, handler HandlerFunc) {
	g.addRoute(http.MethodGet, path, handler)
}

func (g *RouterGroup) POST(path string, handler HandlerFunc) {
	g.addRoute(http.MethodPost, path, handler)
}

func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}
