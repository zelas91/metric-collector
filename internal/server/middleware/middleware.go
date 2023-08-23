package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/logger"
	"time"
)

var log = logger.New()

func WithLogging(c *gin.Context) {
	start := time.Now()

	c.Next()

	duration := time.Since(start)

	log.Infoln(
		"uri", c.Request.RequestURI,
		"method", c.Request.Method,
		"status", c.Writer.Status(),
		"duration", duration,
		"size", c.Writer.Size(),
		"Content-Type", c.GetHeader("Content-Type"),
	)
}

func SetContextPlain(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	c.Next()
}

func SetContextHTML(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.Next()
}
