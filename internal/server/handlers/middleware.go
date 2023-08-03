package handlers

import "github.com/gin-gonic/gin"

func MiddlewareSetContextPlain(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	c.Next()
}

func MiddlewareSetContextHtml(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.Next()
}
