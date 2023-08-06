package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	api := gin.Default()

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api.Run(":8080") // listen and serve on localhost:8080
}
