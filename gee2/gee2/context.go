package gee2

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type HandlerFunc func(c *Context)

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request

	Path   string
	Method string

	Params map[string]string

	handlers []HandlerFunc
	index    int
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		Params: make(map[string]string),
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) HTML(code int, data interface{}) {
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	//c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(data.(string)))
}

func (c *Context) String(code int, format string, args ...interface{}) {
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.WriteHeader(code) // 先设置状态码
	c.Writer.Write([]byte(fmt.Sprintf(format, args...)))
}

func (c *Context) Fail(code int, err string) {
	c.JSON(code, H{"error": err})
}

func (c *Context) JSON(code int, data H) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(data); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) Abort() {
	c.index = len(c.handlers)
}
