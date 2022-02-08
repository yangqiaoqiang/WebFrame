package main

import (
	"net/http"

	"gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello geektutu\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *gin.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
