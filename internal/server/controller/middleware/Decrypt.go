package middleware

import (
	"bytes"
	"crypto/rsa"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/utils/crypto"
	"io"
)

func Decrypt(key *rsa.PrivateKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		if key != nil {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				log.Errorf("Failed to read request body:%v", err)
				c.AbortWithStatus(500)
				return
			}
			defer func() {
				if err := c.Request.Body.Close(); err != nil {
					log.Errorf("body close err:%v", err)
				}
			}()
			body, err = crypto.Decrypt(key, body)
			if err != nil {
				log.Errorf("Failed to decrypt request body:%v", err)
				c.AbortWithStatus(400)
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
		}
		c.Next()
	}
}
