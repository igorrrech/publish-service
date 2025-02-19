package readpost

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jwtauth "github.com/igorrrech/publish-service/publications/pkg/jwtAuth"
	h "github.com/igorrrech/publish-service/publications/service/handlers"
	"github.com/igorrrech/publish-service/publications/service/models"
)

type Request struct {
	GroupID uint `json:"groupId"`
}
type Response struct {
	Posts []models.Post `json:"posts"`
}
type PostReader interface {
	ReadAllPostsInGroup(groupId uint) ([]models.Post, error)
}
type UserGetter interface {
	GetUserById(user_id uint) (models.User, error)
}

func ReadAll(
	pr PostReader,
	ug UserGetter,
) gin.HandlerFunc {
	allowedRols := make(map[string]struct{})
	allowedRols[jwtauth.UserRole] = struct{}{}
	allowedRols[jwtauth.ManagerRole] = struct{}{}
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
		//check user in db
		user, err := ug.GetUserById(jwtUser.UUID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": h.ErrUnknownUser})
			return
		}
		//if user have not access to the group at the hierarchy
		userAllowedGroups := user.GetAllParentGroups()
		if !userAllowedGroups[body.GroupID] {
			c.JSON(http.StatusNotFound, gin.H{"error": h.ErrUserHaveNotAccessToTheGroup})
			return
		}
		//read
		readed, err := pr.ReadAllPostsInGroup(body.GroupID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, Response{
			Posts: readed,
		})
	}
}
