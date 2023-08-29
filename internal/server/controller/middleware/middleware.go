package middleware

import (
	"compress/zlib"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/logger"
	"io"
	"strings"
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

func GzipMiddleware(c *gin.Context) {
	if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "deflate") {
		c.Next()
		return
	}

	writer, _ := zlib.NewWriterLevel(c.Writer, zlib.BestCompression)

	defer func(writer *zlib.Writer) {
		if err := writer.Close(); err != nil {
			log.Debugf("close gzip writer error=%v", err)
		}
	}(writer)

	c.Writer = &zlibWriter{writer, c.Writer}
	c.Header("Content-Encoding", "deflate")

	c.Next()
}

type zlibWriter struct {
	writer *zlib.Writer
	gin.ResponseWriter
}

func (gw *zlibWriter) Write(data []byte) (int, error) {
	return gw.writer.Write(data)
}

func (gw *zlibWriter) WriteString(s string) (int, error) {
	return io.WriteString(gw.writer, s)
}
