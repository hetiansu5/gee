package main

import (
	"fmt"
	"github.com/hetiansu5/gee"
	"net/http"
)

func main() {
	engine := gee.New()
	engine.Get("/my/", func(c *gee.Context) {
		fmt.Println("fire my request")
		c.JSON(http.StatusOK, gee.H{
			"username": "1111",
			"password": "2222",
		})
	})
	engine.Get("/my/:name", func(c *gee.Context) {
		fmt.Println("fire my :name request")
		c.HTML(http.StatusOK, "<p>Hi~</p>"+c.Param("name"))
	})
	v1 := engine.Group("/v1")
	v1.Use(middleware1)

	v2 := v1.Group("/work")
	v2.Use(middleware2)

	v3 := v1.Group("/desk")
	v3.Use(middleware2, middleware1)

	v2.Get("fire", func(c *gee.Context) {
		fmt.Println("fire work request")
		c.Data(http.StatusOK, []byte("work fire"))
	})

	v3.Get("*id/done", func(c *gee.Context) {
		fmt.Println("fire done request")
		c.Data(http.StatusOK, []byte("desk done "+c.Param("id")))
	})
	v3.Get(":name/fine", func(c *gee.Context) {
		fmt.Println("fire fine request")
		c.Data(http.StatusOK, []byte("desk fine "+c.Param("name")))
	})
	engine.Run(":9090")
}

func middleware1(c *gee.Context) {
	fmt.Println("start middleware1")
	c.Next()
	fmt.Println("end middleware1")
}

func middleware2(c *gee.Context) {
	fmt.Println("start middleware2")
	fmt.Println("end middleware2")
}
