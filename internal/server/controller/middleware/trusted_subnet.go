package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zelas91/metric-collector/internal/server/payload"
	"net"
	"net/http"
)

func TrustedSubnet(subnet *net.IPNet) gin.HandlerFunc {
	return func(c *gin.Context) {
		agentIP := c.Request.Header.Get("X-Real-IP")
		if agentIP != "" {
			parseIP := net.ParseIP(agentIP)
			if parseIP == nil || !subnet.Contains(parseIP) {
				payload.NewErrorResponse(c, http.StatusForbidden, "")
				return
			}
		}
		c.Next()
	}
}
