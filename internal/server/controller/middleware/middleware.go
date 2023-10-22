package middleware

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/payload"
	"net/http"
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
