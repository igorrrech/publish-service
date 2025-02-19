package deletepost

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jwtauth "github.com/igorrrech/publish-service/publications/pkg/jwtAuth"
	h "github.com/igorrrech/publish-service/publications/service/handlers"
)

type Request struct {
}

type Response struct {
}
type PostDeleter interface {
	DeletePost() error
}

func Delete(pd PostDeleter) gin.HandlerFunc {
	allowedRols := make(map[string]struct{})
	allowedRols[jwtauth.AdminRole] = struct{}{}

	return func(c *gin.Context) {
		var body Request
		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//get user from context
		value, ok := c.Get("user")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": h.ErrUnknownUser.Error()})
			return
		}
		jwtUser, ok := value.(*jwtauth.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": h.ErrCanNotParseUser.Error()})
			return
		}
		//check premissions
		if _, ok := allowedRols[jwtUser.Role]; !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": h.ErrNotAllowedByRol.Error()})
			return
		}
		//delete
		if err := pd.DeletePost(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, Response{})
	}
}
