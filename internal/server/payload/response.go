package payload

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/logger"
)

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	if err := c.AbortWithError(statusCode, errors.New(message)); err != nil {
		logger.Log.Debugf("Error request status code = %d , error=%v", statusCode, err)
	}

}
