package payload

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	if err := c.AbortWithError(statusCode, errors.New(message)); err != nil {
		logrus.Debugf("Error request status code = %d , error=%v", statusCode, err)
	}

}
