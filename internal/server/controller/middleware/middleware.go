package middleware

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/payload"
	"github.com/zelas91/metric-collector/internal/utils"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	log = logger.New()

	gzipWritePool = &sync.Pool{
		New: func() interface{} {
			return gzip.NewWriter(nil)
		},
	}
)

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

type calculateWriterHash struct {
	gin.ResponseWriter
	body []byte
	key  string
}

func generateHash(cw *calculateWriterHash) error {
	hash, err := utils.GenerateHash(cw.body, cw.key)

	if err != nil {
		if !errors.Is(err, utils.ErrInvalidKey) {
			return fmt.Errorf("calculate hash genetate hash err:%w", err)
		}
		log.Errorf("Invalid hash key")
	}

	if hash != nil {
		cw.Header().Set("HashSHA256", *hash)
	}
	return nil
}

// Write implementation
func (cw *calculateWriterHash) Write(b []byte) (int, error) {
	cw.body = append(cw.body, b...)
	if err := generateHash(cw); err != nil {
		return -1, err
	}
	return cw.ResponseWriter.Write(b)
}

// WriteString implementation
func (cw *calculateWriterHash) WriteString(b string) (int, error) {
	cw.body = append(cw.body, b...)

	if err := generateHash(cw); err != nil {
		return -1, err
	}
	return cw.ResponseWriter.WriteString(b)
}

func CalculateHash(key *string) gin.HandlerFunc {
	return func(c *gin.Context) {

		if key == nil || len(*key) <= 0 {
			c.Next()
			return
		}

		calcWriter := &calculateWriterHash{ResponseWriter: c.Writer, key: *key}
		c.Writer = calcWriter
		c.Next()
	}
}

func HashCheck(key *string) gin.HandlerFunc {
	return func(c *gin.Context) {

		if len(c.GetHeader("HashSHA256")) <= 0 {
			c.Next()
			return
		}
		if key == nil || len(*key) <= 0 {
			c.Next()
			return
		}
		hashKey, err := base64.StdEncoding.DecodeString(*key)
		if err != nil {
			log.Errorf("hash check decode key err:%v", err)
			payload.NewErrorResponseJSON(c, http.StatusInternalServerError, "hash check decode err")
			return
		}
		h := hmac.New(sha256.New, hashKey)

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Errorf("hash check read body err:%v", err)
			payload.NewErrorResponseJSON(c, http.StatusInternalServerError, "hash check read body")
			return
		}
		defer func() {
			if err := c.Request.Body.Close(); err != nil {
				log.Errorf("new check body close err: %v", err)
			}
		}()

		if _, err = h.Write(body); err != nil {
			log.Errorf("hash check generate hash err:%v", err)
			payload.NewErrorResponseJSON(c, http.StatusInternalServerError, "hash check generate hash err")
			return
		}

		hash, err := base64.StdEncoding.DecodeString(c.GetHeader("HashSHA256"))
		if err != nil {
			log.Errorf("hash check decode header hashSHA256 err:%v", err)
			payload.NewErrorResponseJSON(c, http.StatusInternalServerError, "hash check decode header hashSHA256 err")
			return
		}

		if !hmac.Equal(hash, h.Sum(nil)) {
			payload.NewErrorResponse(c, http.StatusBadRequest, "not equal hash")
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		c.Next()
	}
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
			payload.NewErrorResponse(c, http.StatusGatewayTimeout, err.Error())
			return
		}
		payload.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
}
