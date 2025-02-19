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
	userModuleString = "user repository"
)

type UserRepositrory struct {
	db                *gorm.DB
	logger            *slog.Logger
	operationsTimeout time.Duration
	ctx               context.Context
}

func NewUserRepository(
	dsn string,
	logger *slog.Logger,
	operationTimeout time.Duration,
) *UserRepositrory {
	ctx := context.Background()
	db, err := dbgorm.ConnectToDbPg(ctx, dsn, 5*time.Second)
	if err != nil {
		logger.Error(userModuleString, "connection error", err.Error())
	}
	db.AutoMigrate(&models.Group{})
	db.AutoMigrate(&models.User{})
	return &UserRepositrory{
		db:                db,
		logger:            logger,
		operationsTimeout: operationTimeout,
		ctx:               ctx,
	}
}
func (r UserRepositrory) GetUserById(user_id uint) (models.User, error) {
	selectCtx, selectCancel := context.WithTimeout(r.ctx, r.operationsTimeout)
	defer selectCancel()
	var finded models.User
	err := r.db.WithContext(selectCtx).Where("id = ?", user_id).Find(&finded).Error
	return finded, err
}
