package login

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/igorrrech/publish-service/authorization/service/models"
)

const (
	URL = "login"
)

type Request struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
type Response struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
type UserGetter interface {
	GetUserByPhone(string) (*models.User, error)
}
type TokenMaker interface {
	MakeTokenPair(u models.User) (models.TokenPair, error)
}

func Login(
	accessTTL time.Duration,
	refreshTTL time.Duration,
	users UserGetter,
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
		//запрос к репо юзеров
		user, err := users.GetUserByPhone(body.Phone)
		//если пользователя нет возврат ошибки
		//если пароль не правильный возврат ошибки
		if err != nil || user.Password != body.Password {
			c.JSON(http.StatusUnauthorized,
				gin.H{
					"error": err.Error(),
				})
			return
		}
		//создание пары токенов в репо (рефреш - токен)
		pair, err := tokens.MakeTokenPair(*user)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{
					"error": err.Error(),
				})
			return
		}
		//возврат токена с данными о пользователе(ид, роль)
		c.JSON(http.StatusOK, &Response{
			Access:  string(pair.Access),
			Refresh: string(pair.Refresh),
		})
	}
}
