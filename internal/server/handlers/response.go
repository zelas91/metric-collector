package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	if err := c.AbortWithError(statusCode, errors.New(message)); err != nil {

	}

}
