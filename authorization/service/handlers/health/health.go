package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	URL = "health"
)

func Health() gin.HandlerFunc {
	OkMessage := "server is alive"
	return func(c *gin.Context) {
		msg, ok := HealthCheck()
		if !ok {
			c.JSON(http.StatusInternalServerError, msg)
			return
		}
		c.JSON(http.StatusOK, OkMessage)
	}
}
func HealthCheck() (string, bool) {
	return "", true
}
