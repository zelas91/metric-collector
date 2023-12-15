package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/server/payload"
	"github.com/zelas91/metric-collector/internal/utils"
	"io"
	"net/http"
)

type calculateWriterHash struct {
	gin.ResponseWriter
	body *bytes.Buffer
	key  string
}

// Write implementation.
func (cw *calculateWriterHash) Write(b []byte) (int, error) {
	return cw.body.Write(b)
}

// WriteString implementation.
func (cw *calculateWriterHash) WriteString(b string) (int, error) {
	return cw.body.WriteString(b)
}

// CalculateHash middleware.
func CalculateHash(key *string) gin.HandlerFunc {
	return func(c *gin.Context) {

		if key == nil || len(*key) <= 0 {
			c.Next()
			return
		}

		calcWriter := &calculateWriterHash{ResponseWriter: c.Writer, key: *key, body: &bytes.Buffer{}}
		c.Writer = calcWriter
		c.Next()
		/*
			так как после первого вызова ResponseWriter.Write, Written() будет равен true после чего не возможно выставить Заголовок,
			поэтому запись данных из calculateWriterHash.body выполняем уже после вызова next
		*/
		hash, err := utils.GenerateHash(calcWriter.body.Bytes(), calcWriter.key)
		if err != nil {
			if !errors.Is(err, utils.ErrInvalidKey) {
				log.Errorf("calculate hash genetate hash err:%v", err)
				payload.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
			}
			log.Errorf("Invalid hash key")
		}

		if hash != nil {
			c.Header("HashSHA256", *hash)
		}

		if _, err = calcWriter.ResponseWriter.Write(calcWriter.body.Bytes()); err != nil {
			log.Errorf("calculate hash genetate hash err:%v", err)
			payload.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
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
