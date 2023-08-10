package middleware

import "github.com/gin-gonic/gin"

func SetContextPlain(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	c.Next()
}

func SetContextHTML(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.Next()
}
