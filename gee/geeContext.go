package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func Logger() HandlerFunc {
	return func(c *GeeContext) {
		// Start timer
		//t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		//log.Printf("[%d] %s in %v", c.StatusCode, c.r.RequestURI, time.Since(t))
		log.Printf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), c.r.RequestURI)
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
	engine   *Engines
	aborted  bool // 标记是否终止
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

// Abort 终止后续处理
func (c *GeeContext) Abort() {
	c.aborted = true
}

// IsAborted 返回是否已终止
func (c *GeeContext) IsAborted() bool {
	return c.aborted
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
	for c.index < len(c.handlers) && !c.aborted { // 如果未终止，继续执行
		handler := c.handlers[c.index]
		handler(c) // 执行下一个中间件或路由处理函数
		c.index++
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

func (c *GeeContext) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatusCode(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.w, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}

func (c *GeeContext) Fail(code int, msg string) {
	c.SetStatusCode(code)
	c.w.Write([]byte(msg))
	//c.Abort() // 终止后续处理
}

func (c *GeeContext) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatusCode(code)
	encoder := json.NewEncoder(c.w)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.w, err.Error(), 500)
	}
}

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *GeeContext) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
