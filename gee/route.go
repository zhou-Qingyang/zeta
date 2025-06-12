package gee

import (
	"fmt"
	"net/http"
	"strings"
)

type routers struct {
	handlers map[string]HandlerFunc
	roots    map[string]*node
}

func newRouters() *routers {
	return &routers{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	// /p/*
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

func (r *routers) addRoute(method, path string, handler HandlerFunc) {
	//if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" {
	//	panic("method not allowed")
	//}
	//key := method + "-" + path
	//if _, ok := r.handlers[key]; ok {
	//	panic("route already exists")
	//}
	//r.handlers[key] = handler
	parts := parsePattern(path)

	key := method + "-" + path
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(path, parts, 0)
	r.handlers[key] = handler
}

func (r *routers) handle(c *GeeContext) {
	//key := c.Method + "-" + c.Path
	//if handler, ok := r.handlers[key]; ok {
	//	handler(c)
	//} else {
	//	c.String(http.StatusNotFound, "404 page not found")
	//}
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, fmt.Sprintf("404 page not found: %s", c.Path))
	}
}

func (r *routers) getRoute(method string, path string) (*node, map[string]string) {
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
