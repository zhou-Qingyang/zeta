package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type H map[string]any

// type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(c *GeeContext)

type Engines struct {
	*RouterGroup
	router        *routers
	groups        []*RouterGroup     // store all groups
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
}

func (engine *Engines) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engines) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
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

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *GeeContext) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.SetStatusCode(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.w, c.r)
	}
}
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
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
func Default() *Engines {
	engine := New()
	engine.Use(Recovery()) // 确保它是第一个注册的中间件
	//engine.Use(Logger())   // 其他中间件
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
	c.engine = e
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
