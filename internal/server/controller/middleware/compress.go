package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
)

var gzipWritePool = &sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(nil)
	},
}

type gzipWriter struct {
	writer *gzip.Writer
	gin.ResponseWriter
}

// GzipCompress middleware.
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

// GzipDecompress middleware.
func GzipDecompress(c *gin.Context) {
	if c.Request.Header.Get("Content-Encoding") == "gzip" {
		body, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			log.Errorf("gzip decompress new reader err: %v", err)
		}
		defer func() {
			if err := body.Close(); err != nil {
				log.Errorf("GZIP DECOMPRESS BODY CLOSE ERR:%v", err)
			}
		}()

		c.Request.Body = body
		c.Request.Header.Del("Content-Encoding")
		c.Request.Header.Del("Content-Length")
	}
	c.Next()
}
