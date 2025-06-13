package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Logger() HandlerFunc {
	return func(c *GeeContext) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.r.RequestURI, time.Since(t))
	}
}

type GeeContext struct {
	w http.ResponseWriter
	r *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int

	handlers []HandlerFunc
	index    int
}

func NewContext(w http.ResponseWriter, r *http.Request) *GeeContext {
	return &GeeContext{
		w:      w,
		r:      r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1,
	}
}
func (c *GeeContext) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *GeeContext) SetStatusCode(code int) {
	c.StatusCode = code
	c.w.WriteHeader(code)
}

func (c *GeeContext) PostForm(key string) string {
	return c.r.FormValue(key)
}

func (c *GeeContext) Next() {
	c.index++
	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (c *GeeContext) Query(key string) string {
	return c.r.URL.Query().Get(key)
}

func (c *GeeContext) SetHeader(key string, value string) {
	c.w.Header().Set(key, value)
}

func (c *GeeContext) String(code int, value ...string) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatusCode(code)
	////c.w.Write([]byte(value))
	//c.w.Write([]byte(value[0]))
	s := fmt.Sprintf(value[0], value[1:])
	c.w.Write([]byte(s))
}

func (c *GeeContext) Data(code int, data []byte) {
	c.SetStatusCode(code)
	c.w.Write(data)
}

func (c *GeeContext) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatusCode(code)
	c.w.Write([]byte(html))
}

func (c *GeeContext) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatusCode(code)
	encoder := json.NewEncoder(c.w)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.w, err.Error(), 500)
	}
}
