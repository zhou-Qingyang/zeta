package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type H map[string]any

// type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(c *GeeContext)

type Engines struct {
	*RouterGroup
	router *routers
	groups []*RouterGroup // store all groups
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engines      // all groups share a Engine instance
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func New() *Engines {
	engine := &Engines{
		router: newRouters(),
	}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (e *Engines) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := NewContext(w, r)
	c.handlers = middlewares
	e.router.handle(c)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (e *Engines) addRouter(method, path string, handler HandlerFunc) {
	e.router.addRoute(method, path, handler)
}

func (e *Engines) GET(path string, handler HandlerFunc) {
	e.addRouter("GET", path, handler)
}

func (e *Engines) POST(path string, handler HandlerFunc) {
	e.addRouter("POST", path, handler)
}

func (e *Engines) Run(addr string) {
	fmt.Println("ListenAndServe: ", addr)
	log.Fatal(http.ListenAndServe(addr, e))
}
