package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck() gin.HandlerFunc {
	okMessage := "i'm alive"
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, okMessage)
	}
}
