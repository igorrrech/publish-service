package refresh

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/igorrrech/publish-service/authorization/service/models"
)

const (
	URL = "refresh"
)

type Request struct {
	Refresh string `json:"refresh"`
}
type Response struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
type TokenMaker interface {
	GetAccessByRefresh(refresh models.RefreshToken) (*models.TokenPair, error)
}

func Refresh(
	accessTTL time.Duration,
	refreshTTL time.Duration,
	tokens TokenMaker,
) gin.HandlerFunc {

	return func(c *gin.Context) {
		var body Request
		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		//поиск токена по рефрешу
		pair, err := tokens.GetAccessByRefresh(models.RefreshToken(body.Refresh))
		//если не найден, то возврат ошибки
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}
		//возврат токенов
		c.JSON(http.StatusOK, &Response{
			Access:  string(pair.Access),
			Refresh: string(pair.Refresh),
		})
	}
}
