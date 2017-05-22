package main

import (
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") //TODO change
		c.Header("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE")
		c.Header("Access-Control-Allow-Headers", "accept,x-access-token,content-type")

		c.Next()
	}
}
