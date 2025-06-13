package main

import (
	"log"
	"net/http"
	"time"
)

func onlyForV2() HandlerFunc {
	return func(c *GeeContext) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		//c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.r.RequestURI, time.Since(t))
	}
}

func main() {
	//r := NewEngines()
	//r.GET("/", func(c *GeeContext) {
	//	c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	//})
	//r.GET("/hello", func(c *GeeContext) {
	//	// expect /hello?name=geektutu
	//	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	//})
	//r.POST("/login", func(c *GeeContext) {
	//	c.JSON(http.StatusOK, H{
	//		"username": c.PostForm("username"),
	//		"password": c.PostForm("password"),
	//	})
	//})
	//r.Run(":9999")

	r := New()
	r.Use(Logger())
	r.GET("/index", func(c *GeeContext) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *GeeContext) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.GET("/hello", func(c *GeeContext) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	//v2 := r.Group("/v2")
	//{
	//	v2.GET("/hello/:name", func(c *GeeContext) {
	//		// expect /hello/geektutu
	//		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	//	})
	//	v2.POST("/login", func(c *GeeContext) {
	//		c.JSON(http.StatusOK, H{
	//			"username": c.PostForm("username"),
	//			"password": c.PostForm("password"),
	//		})
	//	})
	//}
	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *GeeContext) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}
	r.Run(":9999")
}
