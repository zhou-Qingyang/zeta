package gee2

import (
	"net/http"
	"strings"
)

type Router struct {
	routes map[string]HandlerFunc
	roots  map[string]*node
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]HandlerFunc),
		roots:  make(map[string]*node),
	}
}
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *Router) addRoute(method, pattern string, handlerFunc HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.routes[key] = handlerFunc
}

func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *Router) handle(c *Context) {
	//key := c.Method + ":" + c.Path
	//if handlerFunc, ok := r.routes[key]; ok {
	//	handlerFunc(c)
	//} else {
	//	c.String(404, "404 page not found")
	//}
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.routes[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
	c.Next()
}
