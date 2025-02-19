package repo

import (
	"context"
	"log/slog"
	"time"

	dbgorm "github.com/igorrrech/publish-service/authorization/pkg/dbGorm"
	"github.com/igorrrech/publish-service/publications/service/models"
	"gorm.io/gorm"
)

const (
	postModuleString = "post repository"
)

type PostRepository struct {
	db                *gorm.DB
	logger            *slog.Logger
	operationsTimeout time.Duration
	ctx               context.Context
}

func NewPostRepository(
	dsn string,
	logger *slog.Logger,
	opertionsTimeout time.Duration,
) *PostRepository {
	ctx := context.Background()
	db, err := dbgorm.ConnectToDbPg(ctx, dsn, 5*time.Second)
	if err != nil {
		logger.Error(postModuleString, " connection error", err.Error())
	}
	db.AutoMigrate(&models.Group{})
	db.AutoMigrate(&models.Post{})
	return &PostRepository{
		db:                db,
		logger:            logger,
		operationsTimeout: opertionsTimeout,
		ctx:               ctx,
	}
}
func (r PostRepository) CreatePost(rawpost models.RawPost) (uint, error) {
	var post models.Post
	post.GroupID = rawpost.GroupID
	post.Title = rawpost.Title
	post.Content = rawpost.Content
	if err := post.Validate(); err != nil {
		return 0, err
	}

	createCtx, createCancel := context.WithTimeout(r.ctx, r.operationsTimeout)
	defer createCancel()
	err := r.db.WithContext(createCtx).Create(post).Error
	if err != nil {
		return 0, err
	}
	createCancel()

	selectCtx, selectCancel := context.WithTimeout(r.ctx, r.operationsTimeout)
	defer selectCancel()
	var createdPost models.Post
	err = r.db.Model(&models.Post{}).WithContext(selectCtx).First(&createdPost).Error

	return createdPost.ID, err
}
func (r PostRepository) ReadAllPostsInGroup(groupId uint) ([]models.Post, error) {
	selectCtx, selectCancel := context.WithTimeout(r.ctx, r.operationsTimeout)
	defer selectCancel()
	var posts []models.Post
	err := r.db.WithContext(selectCtx).Where("group_id = ? ", groupId).Find(&posts).Error

	return posts, err
}
func (r PostRepository) UpdatePost() error {
	return nil
}
func (r PostRepository) DeletePost() error {
	return nil
}

// func (r PostRepository) validatePost(post models.Post) error {
// 	if err := post.Validate(); err != nil {
// 		return err
// 	}

// 	groupToFind := models.Group{}
// 	groupToFind.ID = post.GroupID

// 	selectCtx, selectCancel := context.WithTimeout(r.ctx, r.operationsTimeout)
// 	defer selectCancel()
// 	err := r.db.Model(&groupToFind).WithContext(selectCtx).First(&groupToFind).Error

// 	return err
// }
