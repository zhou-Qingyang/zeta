package main

import (
	"github.com/zhou-Qingyang/hrm/ZetaCache/gee"
	"net/http"
)

func main() {
	geeEngine := gee.NewEngines()
	//geeEngine.GET("/", func(c *gee.GeeContext) {
	//	c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	//})
	//geeEngine.GET("/hello", func(c *gee.GeeContext) {
	//	fmt.Println(c.Query("name"))
	//	c.String(http.StatusOK, "hello,you're at")
	//})
	//geeEngine.POST("/login", func(c *gee.GeeContext) {
	//	c.JSON(http.StatusOK, gee.H{
	//		"username": c.PostForm("username"),
	//		"password": c.PostForm("password"),
	//	})
	//})
	geeEngine.GET("/", func(c *gee.GeeContext) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	geeEngine.GET("/hello", func(c *gee.GeeContext) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	geeEngine.GET("/hello/:name", func(c *gee.GeeContext) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	geeEngine.GET("/hello/b/c", func(c *gee.GeeContext) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", "b/c", c.Path)
	})

	geeEngine.GET("/assets/*filepath", func(c *gee.GeeContext) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})

	geeEngine.Run(":9999")
}
