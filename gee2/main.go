package main

import (
	"github.com/zhou-Qingyang/hrm/ZetaCache/gee2/gee2"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gee2.HandlerFunc {
	return func(c *gee2.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("onlyForV2 %s in %v for group v2", c.Req.RequestURI, time.Since(t))
	}
}

func Logger() gee2.HandlerFunc {
	return func(c *gee2.Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("Logger %s in %v", c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	//fmt.Println("Starting Gee2...")
	//http.ListenAndServe(":9999", &GeeContext{})
	//r := gee2.New()
	//r.GET("/", func(c *gee2.Context) {
	//	c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	//})
	//r.GET("/hello", func(c *gee2.Context) {
	//	// expect /hello?name=geektutu
	//	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	//})
	//
	//r.GET("/hello/:name", func(c *gee2.Context) {
	//	// expect /hello/geektutu
	//	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	//})

	//Invoke-WebRequest -Uri "http://localhost:9999/login" -Method POST -Body 'username=geektutu&password=1234' -ContentType "application/x-www-form-urlencoded"
	//r.Run(":9999")

	//r := gee2.New()
	//v1 := r.Group("/v1")
	//v1.GET("/", func(c *gee2.Context) {
	//	c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	//})
	//
	//v2 := r.Group("/v2")
	//v2.GET("/hello", func(c *gee2.Context) {
	//	// expect /hello?name=geektutu
	//	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	//})
	//
	//r.Run(":9999")

	//r := gee2.New()
	//r.Use(Logger()) // global midlleware
	//r.GET("/", func(c *gee2.Context) {
	//	c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	//})
	//v2 := r.Group("/v2")
	//v2.Use(onlyForV2()) // v2 group middleware
	//{
	//	v2.GET("/hello/:name", func(c *gee2.Context) {
	//		// expect /hello/geektutu
	//		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	//	})
	//}
	//r.Run(":9999")

	r := gee2.New()
	r.Use(gee2.Recovery())
	r.GET("/", func(c *gee2.Context) {
		c.String(http.StatusOK, "Hello Geektutu\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *gee2.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})
	r.Run(":9999")
}
