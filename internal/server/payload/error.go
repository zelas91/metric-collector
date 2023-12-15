// Package payload error handler with response capability to the client

package payload

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/logger"
)

var log = logger.New()

// NewErrorResponse response err in string.
func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	if err := c.AbortWithError(statusCode, errors.New(message)); err != nil {
		log.Errorf("Error request status code = %d , error=%v", statusCode, err)
	}

}

// NewErrorResponseJSON response err in JSON.
func NewErrorResponseJSON(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, gin.H{"message": message})
}
