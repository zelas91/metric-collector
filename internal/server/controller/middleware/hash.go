package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/server/payload"
	"github.com/zelas91/metric-collector/internal/utils"
	"io"
	"net/http"
)

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
