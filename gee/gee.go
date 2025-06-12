package gee

import (
	"fmt"
	"log"
	"net/http"
)

type H map[string]any

// type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(c *GeeContext)

type Engines struct {
	routers *routers
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine instance
}

func NewEngines() *Engines {
	return &Engines{
		routers: newRouters(),
	}
}

func (e *Engines) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r)
	e.routers.handle(c)
}

func (e *Engines) addRouter(method, path string, handler HandlerFunc) {
	e.routers.addRoute(method, path, handler)
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
