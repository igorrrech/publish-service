package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	jwtauth "github.com/igorrrech/publish-service/publications/pkg/jwtAuth"
)

type Request struct {
	Access string `json:"access"`
	UserID uint   `json:"userId"`
}

func Auth(
	secret string,
) gin.HandlerFunc {
	s := secret
	return func(c *gin.Context) {
		var body Request
		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, ErrBodyHaveNotClaims)
			return
		}
		user, err := jwtauth.AccessValidate(body.Access, s)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}
		//additional token validate. MITM need to now also uuid in token
		if user.UUID != body.UserID {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrUserNotEquals})
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

var (
	ErrBodyHaveNotClaims = errors.New("have not access token or userId in body")
	ErrUserNotEquals     = errors.New("user_id in body not equals user_id in token")
)
