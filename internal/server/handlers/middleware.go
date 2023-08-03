package handlers

import "github.com/gin-gonic/gin"

func MiddlewareSetContextPlain(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	c.Next()
}

func MiddlewareSetContextHTML(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.Next()
}
