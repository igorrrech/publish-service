package createpost

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jwtauth "github.com/igorrrech/publish-service/publications/pkg/jwtAuth"
	h "github.com/igorrrech/publish-service/publications/service/handlers"
	"github.com/igorrrech/publish-service/publications/service/models"
)

type Request struct {
	models.RawPost
}
type Response struct {
	PostId uint `json:"post-id"`
}
type PostCreator interface {
	CreatePost(rawpost models.RawPost) (uint, error)
}
type UserGetter interface {
	GetUserById(user_id uint) (models.User, error)
}

func Create(
	pc PostCreator,
	ug UserGetter,
) gin.HandlerFunc {
	allowedRols := make(map[string]struct{})
	allowedRols[jwtauth.AdminRole] = struct{}{}
	allowedRols[jwtauth.ManagerRole] = struct{}{}

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
		//check user in db
		user, err := ug.GetUserById(jwtUser.UUID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": h.ErrUnknownUser})
			return
		}
		//check user groups
		if !user.IsInGroup(body.GroupID) {
			c.JSON(http.StatusNotFound, gin.H{"error": h.ErrUserNotInGroup})
			return
		}
		//create
		postId, err := pc.CreatePost(body.RawPost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, Response{
			PostId: postId,
		})
	}
}
