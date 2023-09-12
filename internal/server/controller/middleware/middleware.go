package middleware

import (
	"compress/gzip"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/logger"
	"net/http"
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

func GzipCompress(c *gin.Context) {

	if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
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
func (gw *gzipWriter) WriteString(data string) (int, error) {
	return gw.writer.Write([]byte(data))
}

func getGzipWriter() *gzip.Writer {
	if v := gzipWritePool.Get(); v != nil {
		return v.(*gzip.Writer)
	}
	writer, err := gzip.NewWriterLevel(nil, gzip.BestCompression)
	if err != nil {
		log.Errorf("Failed to create gzip writer err: %v", err)
		return nil
	}

	return writer
}

func releaseGzipWriter(writer *gzip.Writer) {
	defer func(w *gzip.Writer) {
		if err := writer.Close(); err != nil {
			log.Error("Failed to close gzip writer:", err)
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

func Timeout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	c.Request = c.Request.WithContext(ctx)

	ch := make(chan struct{})
	go func() {
		c.Next()
		close(ch)
	}()

	select {
	case <-ch:
		return
	case <-ctx.Done():
		err := ctx.Err()
		if errors.Is(err, context.DeadlineExceeded) {
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{"error": "Request timeout"})
			return
		}
		if err = c.AbortWithError(http.StatusInternalServerError, err); err != nil {
			log.Errorf("timeout middleware err :%v", err)
			return
		}
		return
	}
}
