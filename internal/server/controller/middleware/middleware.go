package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/logger"
	"strings"
	"sync"
	"time"
)

var log = logger.New()
var gzipWritePool = &sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(nil)
	},
}

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

func GzipCompress(c *gin.Context) {

	if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") &&
		(strings.Contains(c.Request.Header.Get("Content-Type"), "application/json") ||
			strings.Contains(c.Request.Header.Get("Content-Type"), "text/html")) {
		c.Next()
		return
	}

	writer := getGzipWriter()
	defer releaseGzipWriter(writer)
	writer.Reset(c.Writer)
	c.Writer = &gzipWriter{writer, c.Writer}

	c.Header("Content-Encoding", "gzip")

	c.Next()
}

type gzipWriter struct {
	writer *gzip.Writer
	gin.ResponseWriter
}

func (gw *gzipWriter) Write(data []byte) (int, error) {

	return gw.writer.Write(data)
}

func getGzipWriter() *gzip.Writer {
	if v := gzipWritePool.Get(); v != nil {
		return v.(*gzip.Writer)
	}
	writer, err := gzip.NewWriterLevel(nil, gzip.BestCompression)
	if err != nil {
		log.Info("Failed to create gzip writer:", err)
		return nil
	}

	return writer
}
func releaseGzipWriter(writer *gzip.Writer) {
	defer func(w *gzip.Writer) {
		if err := writer.Close(); err != nil {
			log.Debug("Failed to close gzip writer:", err)
			return
		}
	}(writer)
	gzipWritePool.Put(writer)
}
func GzipDecompress(c *gin.Context) {
	if c.Request.Header.Get("Content-Encoding") == "gzip" {
		c.Request.Body, _ = gzip.NewReader(c.Request.Body)
		c.Request.Header.Del("Content-Encoding")
		c.Request.Header.Del("Content-Length")
	}
	c.Next()
}
